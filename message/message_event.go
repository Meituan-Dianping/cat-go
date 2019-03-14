package message

import (
	"sync/atomic"
)

type Event struct {
	Message
}

func (e *Event) Complete() {
	if !atomic.CompareAndSwapUint32(&e.isCompleted, 0, 1) {
		return
	}

	if e.Message.flush != nil {
		e.Message.flush(e)
	}
}

func NewEvent(mtype, name string, flush Flush) *Event {
	return &Event{
		Message: NewMessage(mtype, name, flush),
	}
}
