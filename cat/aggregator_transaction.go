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

func (p *transactionAggregator) GetName() string {
	return "TransactionAggregator"
}

type transactionAggregator struct {
	signalsMixin
	ch      chan *message.Transaction
	dataMap map[string]*transactionData
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

func (p *transactionAggregator) BackGround() {
	var ticker = time.NewTicker(transactionAggregatorInterval)
	for p.isAlive {
		select {
		case signal := <-p.signals:
			if signal == signalShutdown {
				close(p.ch)
				ticker.Stop()
				p.stop()
			}
		case trans := <-p.ch:
			p.getOrDefault(trans).add(trans)
		case <-ticker.C:
			p.collectAndSend()
		}
	}

	p.collectAndSend()
	p.exit()
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

	if transaction.GetStatus() != SUCCESS {
		data.fail++
	}

	millis := transaction.GetDuration().Nanoseconds() / time.Millisecond.Nanoseconds()
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
		signalsMixin: makeSignalsMixedIn(signalTransactionAggregatorExit),
		ch:           make(chan *message.Transaction, transactionAggregatorChannelCapacity),
		dataMap:      make(map[string]*transactionData),
	}
}
