package message

type Heartbeat struct {
	Message
}

func (e *Heartbeat) Complete() {
	if e.Message.flush != nil {
		e.Message.flush(e)
	}
}

func NewHeartbeat(mtype, name string, flush Flush) *Heartbeat {
	return &Heartbeat{
		Message: NewMessage(mtype, name, flush),
	}
}
