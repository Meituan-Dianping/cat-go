package cat

import (
	"bytes"
	"strconv"

	"../message"
)

const batchFlag = '@'
const batchSplit = ';'

type LocalAggregator struct {
	event *eventAggregator
	transaction *transactionAggregator
	metric *metricAggregator
}

func (p *LocalAggregator) flush(m message.Messager) {
	switch m := m.(type) {
		case *message.Transaction:
			sender.handleTransaction(m)
		default:
			logger.Warning("Aggregator flusher expected a transaction.")
	}
}

func (p *LocalAggregator) Background() {
	go p.event.BackGround()
	go p.transaction.BackGround()
	go p.metric.BackGround()
}

type Buf struct {
	bytes.Buffer
}

func newBuf() *Buf {
	return &Buf{
		*bytes.NewBuffer([]byte{}),
	}
}

func (b *Buf) WriteInt(i int) {
	b.WriteString(strconv.Itoa(i))
}

func (b *Buf) WriteUInt64(i uint64) {
	b.WriteString(strconv.FormatUint(i, 10))
}

func computeDuration(durationInMillis int) int {
	if durationInMillis < 1 {
		return 1
	} else if durationInMillis < 20 {
		return durationInMillis
	} else if durationInMillis < 200 {
		return durationInMillis - durationInMillis% 5
	} else if durationInMillis < 500 {
		return durationInMillis - durationInMillis% 20
	} else if durationInMillis < 2000 {
		return durationInMillis - durationInMillis% 50
	} else if durationInMillis < 20000 {
		return durationInMillis - durationInMillis% 500
	} else if durationInMillis < 1000000 {
		return durationInMillis - durationInMillis% 10000
	} else {
		dk := 524288
		if durationInMillis > 3600 * 1000 {
			dk = 3600 * 1000
		} else {
			for dk < durationInMillis {
				dk <<= 1
			}
		}
		return dk
	}
}

var aggregator = LocalAggregator{
	event: newEventAggregator(),
	transaction: newTransactionAggregator(),
	metric: newMetricAggregator(),
}