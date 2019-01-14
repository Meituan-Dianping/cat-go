package cat

import (
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

type transactionAggregator struct {
	ch      chan *message.Transaction
	dataMap map[string]*transactionData
}

func (p *transactionAggregator) send(dataMap map[string]*transactionData) {
	if len(dataMap) == 0 {
		return
	}

	t := message.NewTransaction(CAT_SYSTEM, TRANSACTION_AGGREGATOR, aggregator.flush)
	defer t.Complete()

	buf := newBuf()

	for _, data := range dataMap {
		event := t.NewEvent(data.mtype, data.name)
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

		event.SetData(buf.String())
	}
}

func (p *transactionAggregator) getOrDefaultData(transaction *message.Transaction) *transactionData {
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

func (p *transactionAggregator) BackGround() {
	var ticker = time.NewTicker(time.Second)
	for {
		select {
		case trans := <-p.ch:
			p.getOrDefaultData(trans).add(trans)
		case <-ticker.C:
			dataMap := p.dataMap
			p.dataMap = make(map[string]*transactionData)
			p.send(dataMap)
		}
	}
}

func (p *transactionAggregator) Put(t *message.Transaction) {
	select {
	case p.ch <- t:
	default:
		logger.Warning("Transaction aggregator is full")
	}
}

func (data *transactionData) add(transaction *message.Transaction) {
	data.count++

	if transaction.GetStatus() != CAT_SUCCESS {
		data.fail++
	}

	data.sum += transaction.GetDurationInMillis()

	duration := computeDuration(int(transaction.GetDurationInMillis()))
	if _, ok := data.durations[duration]; ok {
		data.durations[duration]++
	} else {
		data.durations[duration] = 1
	}
}

func newTransactionAggregator() *transactionAggregator {
	return &transactionAggregator{
		ch:      make(chan *message.Transaction, TRANSACTION_AGGREGATOR_CHANNEL_CAPACITY),
		dataMap: make(map[string]*transactionData),
	}
}
