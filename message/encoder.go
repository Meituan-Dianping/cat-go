package message

import (
	"bytes"
)

const (
	defaultThreadGroupName = ""
	defaultThreadId        = "0"
	defaultThreadName      = ""
)

type Encoder interface {
	EncodeHeader(*bytes.Buffer, *Header) error
	EncodeMessage(*bytes.Buffer, Messager) error
	EncodeTransaction(*bytes.Buffer, *Transaction) error
	EncodeEvent(*bytes.Buffer, *Event) error
	EncodeHeartbeat(*bytes.Buffer, *Heartbeat) error
	EncodeMetric(*bytes.Buffer, *Metric) error
}

type encoderBase struct {
}

func encodeMessage(encoder Encoder, buf *bytes.Buffer, message Messager) (err error) {
	switch m := message.(type) {
	case *Transaction:
		return encoder.EncodeTransaction(buf, m)
	case *Event:
		return encoder.EncodeEvent(buf, m)
	case *Heartbeat:
		return encoder.EncodeHeartbeat(buf, m)
	case *Metric:
		return encoder.EncodeMetric(buf, m)
	default:
		// TODO unsupported message type.
		return
	}
}


