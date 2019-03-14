package message

import (
	"sync/atomic"
)

type Metric struct {
	Message
}

func (e *Metric) Complete() {
	if !atomic.CompareAndSwapUint32(&e.isCompleted, 0, 1) {
		return
	}

	if e.Message.flush != nil {
		e.Message.flush(e)
	}
}

func NewMetric(mtype, name string, flush Flush) *Metric {
	return &Metric{
		Message: NewMessage(mtype, name, flush),
	}
}
