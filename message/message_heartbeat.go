package message

import (
	"sync/atomic"
)

type Heartbeat struct {
	Message
}

func (e *Heartbeat) Complete() {
	if !atomic.CompareAndSwapUint32(&e.isCompleted, 0, 1) {
		return
	}

	if e.Message.flush != nil {
		e.Message.flush(e)
	}
}

func NewHeartbeat(mtype, name string, flush Flush) *Heartbeat {
	return &Heartbeat{
		Message: NewMessage(mtype, name, flush),
	}
}
