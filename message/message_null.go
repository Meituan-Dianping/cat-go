package message

import (
	"bytes"
	"time"
)

type nullMessage struct {
}

type nullTransaction struct {
	nullMessage
}

var NullMessage = &nullMessage{}

var NullTransaction = &nullTransaction{}

func (m *nullMessage) Complete() {
}

func (m *nullMessage) GetType() string {
	return ""
}

func (m *nullMessage) GetName() string {
	return ""
}

func (m *nullMessage) GetStatus() string {
	return ""
}

func (m *nullMessage) GetData() *bytes.Buffer {
	return nil
}

func (m *nullMessage) GetTime() time.Time {
	return time.Now()
}

func (m *nullMessage) SetTime(time time.Time) {
}

func (m *nullMessage) SetTimestamp(timestampMs int64) {
}

func (m *nullMessage) AddData(k string, v ...string) {
}

func (m *nullMessage) SetData(v string) {
}

func (m *nullMessage) SetStatus(status string) {
}

func (t *nullTransaction) GetChildren() []Messager {
	return []Messager{}
}

func (t *nullTransaction) GetDuration() time.Duration {
	return 0
}

func (t *nullTransaction) SetDuration(duration time.Duration) {
}

func (t *nullTransaction) SetDurationStart(duration time.Time) {
}

func (t *nullTransaction) NewEvent(mtype, name string) Messager {
	return NullMessage
}

func (t *nullTransaction) LogEvent(mtype, name string, args ...string) {
	return
}
