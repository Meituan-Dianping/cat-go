package cat

import (
	"bytes"
	"fmt"
	"time"

	"../message"
)

type transactionData struct {
	mtype, name string

	count, fail int

	sum int64

	durations map[int]int
}

//noinspection GoUnhandledErrorResult
func encodeTransactionData(data *transactionData) *bytes.Buffer {
	buf := newBuf()

	buf.WriteRune(batchFlag)
	buf.WriteInt(data.count)
	buf.WriteRune(batchSplit)
	buf.WriteInt(data.fail)
	buf.WriteRune(batchSplit)
	buf.WriteUInt64(uint64(data.sum))
	buf.WriteRune(batchSplit)

	i := 0
	for k, v := range data.durations {
		if i > 0 {
			buf.WriteRune('|')
		}
		buf.WriteInt(k)
		buf.WriteRune(',')
		buf.WriteInt(v)
		i++
	}
	buf.WriteRune(batchSplit)

	return &buf.Buffer
}

func (p *transactionAggregator) GetName() string {
	return "TransactionAggregator"
}

type transactionAggregator struct {
	scheduleMixin
	ch      chan *message.Transaction
	dataMap map[string]*transactionData
	ticker *time.Ticker
}

func (p *transactionAggregator) collectAndSend() {
	dataMap := p.dataMap
	p.dataMap = make(map[string]*transactionData)
	p.send(dataMap)
}

func (p *transactionAggregator) send(dataMap map[string]*transactionData) {
	if len(dataMap) == 0 {
		return
	}

	t := message.NewTransaction(typeSystem, nameTransactionAggregator, aggregator.flush)
	defer t.Complete()

	for _, data := range dataMap {
		trans := message.NewTransaction(data.mtype, data.name, nil)
		trans.SetData(encodeTransactionData(data).String())
		trans.Complete()
		t.AddChild(trans)
	}
}

func (p *transactionAggregator) getOrDefault(transaction *message.Transaction) *transactionData {
	key := fmt.Sprintf("%s,%s", transaction.Type, transaction.Name)

	if data, ok := p.dataMap[key]; ok {
		return data
	} else {
		p.dataMap[key] = &transactionData{
			mtype:     transaction.GetType(),
			name:      transaction.GetName(),
			count:     0,
			fail:      0,
			sum:       0,
			durations: make(map[int]int),
		}
		return p.dataMap[key]
	}
}

func (p *transactionAggregator) afterStart() {
	p.ticker = time.NewTicker(transactionAggregatorInterval)
}

func (p *transactionAggregator) beforeStop() {
	close(p.ch)

	for t := range p.ch {
		p.getOrDefault(t).add(t)
	}
	p.collectAndSend()

	p.ticker.Stop()
}

func (p *transactionAggregator) process() {
	select {
	case sig := <-p.signals:
		p.handle(sig)
	case t := <-p.ch:
		p.getOrDefault(t).add(t)
	case <-p.ticker.C:
		p.collectAndSend()
	}
}

func (p *transactionAggregator) Put(t *message.Transaction) {
	if !IsEnabled() {
		return
	}

	select {
	case p.ch <- t:
	default:
		logger.Warning("Transaction aggregator is full")
	}
}

func (data *transactionData) add(transaction *message.Transaction) {
	data.count++

	if transaction.GetStatus() != SUCCESS {
		data.fail++
	}

	millis := duration2Millis(transaction.GetDuration())
	data.sum += millis

	duration := computeDuration(int(millis))
	if _, ok := data.durations[duration]; ok {
		data.durations[duration]++
	} else {
		data.durations[duration] = 1
	}
}

func newTransactionAggregator() *transactionAggregator {
	return &transactionAggregator{
		scheduleMixin: makeScheduleMixedIn(signalTransactionAggregatorExit),
		ch:            make(chan *message.Transaction, transactionAggregatorChannelCapacity),
		dataMap:       make(map[string]*transactionData),
	}
}
