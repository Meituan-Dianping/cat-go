package message

type Event struct {
	Message
}

func (e *Event) Complete() {
	if e.Message.flush != nil {
		e.Message.flush(e)
	}
}

func NewEvent(mtype, name string, flush Flush) *Event {
	return &Event{
		Message: NewMessage(mtype, name, flush),
	}
}
