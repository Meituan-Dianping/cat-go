package message

import (
	"sync"
	"time"
)

type TransactionGetter interface {
	GetChildren() []Messager
	GetDuration() time.Duration
}

type Transactor interface {
	Messager
	TransactionGetter
	SetDuration(duration time.Duration)
	SetDurationStart(time time.Time)
	NewEvent(mtype, mname string) Messager
	LogEvent(mtype, mname string, args ...string)
}

type Transaction struct {
	Message

	children []Messager

	isCompleted bool

	mu sync.Mutex

	duration time.Duration
	durationStart time.Time
}

func (t *Transaction) Complete() {
	if t.isCompleted {
		return
	}
	t.isCompleted = true

	if t.duration == 0 {
		t.duration = time.Now().Sub(t.Message.timestamp)
	}

	if t.Message.flush != nil {
		t.Message.flush(t)
	}
}

func (t *Transaction) GetChildren() []Messager {
	return t.children
}

func (t *Transaction) GetDuration() time.Duration {
	return t.duration
}

func (t *Transaction) GetDurationInMillis() int64 {
	return t.duration.Nanoseconds() / time.Millisecond.Nanoseconds()
}

func (t *Transaction) SetDuration(duration time.Duration) {
	t.duration = duration
}
func (t *Transaction) SetDurationStart(time time.Time) {
	t.durationStart = time
}

func (t *Transaction) NewEvent(mtype, mname string) Messager {
	var e = NewEvent(mtype, mname, nil)
	t.AddChild(e)
	return e
}

func (t *Transaction) LogEvent(mtype, mname string, args ...string) {
	var e = t.NewEvent(mtype, mname)
	if len(args) > 0 {
		e.SetStatus(args[0])
	}
	if len(args) > 1 {
		e.AddDataPair(args[1])
	}
	e.Complete()
}

func (t *Transaction) AddChild(messager Messager) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.children = append(t.children, messager)
}

func NewTransaction(mtype, name string, flush Flush) *Transaction {
	return &Transaction{
		Message:       NewMessage(mtype, name, flush),
		children:      make([]Messager, 0),
		isCompleted:   false,
		mu:            sync.Mutex{},
		duration:      0,
		durationStart: time.Time{},
	}
}
