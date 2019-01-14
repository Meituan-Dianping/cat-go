package message

import (
	"bytes"
	"time"
)

const (
	CAT_SUCCESS = "0"
	CAT_ERROR   = "-1"
)

type Flush func(m Messager)

type MessageGetter interface {
	GetType() string
	GetName() string
	GetStatus() string
	GetData() *bytes.Buffer
	GetTime() time.Time
}

type Messager interface {
	MessageGetter
	AddData(k string, v ...string)
	SetData(v string)
	SetStatus(status string)
	SetTime(time time.Time)
	Complete()
}

type Message struct {
	Type   string
	Name   string
	Status string

	timestamp time.Time

	data *bytes.Buffer

	flush Flush
}

func NewMessage(mtype, name string, flush Flush) Message {
	return Message{
		Type:      mtype,
		Name:      name,
		Status:    CAT_SUCCESS,
		timestamp: time.Now(),
		data:      new(bytes.Buffer),
		flush:     flush,
	}
}

func (m *Message) Complete() {
	if m.flush != nil {
		m.flush(m)
	}
}

func (m *Message) GetType() string {
	return m.Type
}

func (m *Message) GetName() string {
	return m.Name
}

func (m *Message) GetStatus() string {
	return m.Status
}

func (m *Message) GetData() *bytes.Buffer {
	return m.data
}

func (m *Message) GetTime() time.Time {
	return m.timestamp
}

func (m *Message) SetTime(t time.Time) {
	m.timestamp = t
}

func (m *Message) AddData(k string, v ...string) {
	if m.data.Len() != 0 {
		m.data.WriteRune('&')
	}
	if len(v) == 0 {
		m.data.WriteString(k)
	} else {
		m.data.WriteString(k)
		m.data.WriteRune('=')
		m.data.WriteString(v[0])
	}
}

func (m *Message) SetData(v string) {
	m.data.Reset()
	m.data.WriteString(v)
}

func (m *Message) SetStatus(status string) {
	m.Status = status
}

func (m *Message) SetSuccessStatus() {
	m.Status = CAT_SUCCESS
}
