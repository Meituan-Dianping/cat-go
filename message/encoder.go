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
}