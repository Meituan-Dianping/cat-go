package message

import (
	"bytes"
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

func (e *encoderBase) EncodeHeader(buf *bytes.Buffer, header *Header) (err error) {
	if _, err = buf.WriteString(BINARY_PROTOCOL); err != nil {
		return
	}
	if err = writeString(buf, header.Domain); err != nil {
		return
	}
	if err = writeString(buf, header.Hostname); err != nil {
		return
	}
	if err = writeString(buf, header.Ip); err != nil {
		return
	}

	if err = writeString(buf, defaultThreadGroupName); err != nil {
		return
	}
	if err = writeString(buf, defaultThreadId); err != nil {
		return
	}
	if err = writeString(buf, defaultThreadName); err != nil {
		return
	}

	if err = writeString(buf, header.MessageId); err != nil {
		return
	}
	if err = writeString(buf, header.ParentMessageId); err != nil {
		return
	}
	if err = writeString(buf, header.RootMessageId); err != nil {
		return
	}

	// sessionToken.
	if err = writeString(buf, ""); err != nil {
		return
	}
	return
}
