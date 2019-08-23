package message

import (
	"bytes"
	"time"
)



type BinaryEncoder struct {
	encoderBase
}

func NewBinaryEncoder() *BinaryEncoder {
	return &BinaryEncoder{}
}

func writeI64(buf *bytes.Buffer, i int64) (err error) {
	for {
		if i&^0x7F == 0 {
			if err = buf.WriteByte(byte(i)); err != nil {
				return
			}
			return
		} else {
			if err = buf.WriteByte(byte(i&0x7F | 0x80)); err != nil {
				return
			}
			i >>= 7
		}
	}
}

func (e *BinaryEncoder) writeString(buf *bytes.Buffer, s string) (err error) {
	if err = writeI64(buf, int64(len(s))); err != nil {
		return
	}
	if _, err = buf.WriteString(s); err != nil {
		return
	}
	return
}

func (e *BinaryEncoder) EncodeHeader(buf *bytes.Buffer, header *Header) (err error) {
	if _, err = buf.WriteString(BinaryProtocol); err != nil {
		return
	}
	if err = e.writeString(buf, header.Domain); err != nil {
		return
	}
	if err = e.writeString(buf, header.Hostname); err != nil {
		return
	}
	if err = e.writeString(buf, header.Ip); err != nil {
		return
	}

	if err = e.writeString(buf, defaultThreadGroupName); err != nil {
		return
	}
	if err = e.writeString(buf, defaultThreadId); err != nil {
		return
	}
	if err = e.writeString(buf, defaultThreadName); err != nil {
		return
	}

	if err = e.writeString(buf, header.MessageId); err != nil {
		return
	}
	if err = e.writeString(buf, header.ParentMessageId); err != nil {
		return
	}
	if err = e.writeString(buf, header.RootMessageId); err != nil {
		return
	}

	// sessionToken.
	if err = e.writeString(buf, ""); err != nil {
		return
	}
	return
}

func (e *BinaryEncoder) encodeMessageStart(buf *bytes.Buffer, m Messager) (err error) {
	var timestamp = m.GetTime().UnixNano() / time.Millisecond.Nanoseconds()
	if err = writeI64(buf, timestamp); err != nil {
		return
	}
	if err = e.writeString(buf, m.GetType()); err != nil {
		return
	}
	if err = e.writeString(buf, m.GetName()); err != nil {
		return
	}
	return
}

func (e *BinaryEncoder) encodeMessageEnd(buf *bytes.Buffer, m Messager) (err error) {
	if err = e.writeString(buf, m.GetStatus()); err != nil {
		return
	}

	if m.GetData() == nil {
		if err = writeI64(buf, 0); err != nil {
			return
		}
	} else {
		if err = writeI64(buf, int64(m.GetData().Len())); err != nil {
			return
		}
		if _, err = buf.Write(m.GetData().Bytes()); err != nil {
			return
		}
	}
	return
}

func (e *BinaryEncoder)encodeMessageWithLeader(buf *bytes.Buffer, m *Message, leader rune) (err error) {
	if _, err = buf.WriteRune(leader); err != nil {
		return
	}
	if err = e.encodeMessageStart(buf, m); err != nil {
		return
	}
	if err = e.encodeMessageEnd(buf, m); err != nil {
		return
	}
	return
}

func (e *BinaryEncoder) EncodeMessage(buf *bytes.Buffer, message Messager) (err error) {
	return encodeMessage(e, buf, message)
}

func (e *BinaryEncoder) EncodeTransaction(buf *bytes.Buffer, trans *Transaction) (err error) {
	if _, err = buf.WriteRune('t'); err != nil {
		return
	}
	if err = e.encodeMessageStart(buf, trans); err != nil {
		return
	}

	for _, message := range trans.GetChildren() {
		if err = e.EncodeMessage(buf, message); err != nil {
			return
		}
	}

	if _, err = buf.WriteRune('T'); err != nil {
		return
	}
	if err = e.encodeMessageEnd(buf, trans); err != nil {
		return
	}

	duration := trans.GetDuration().Nanoseconds() / time.Microsecond.Nanoseconds()
	if err = writeI64(buf, duration); err != nil {
		return
	}
	return
}

func (e *BinaryEncoder) EncodeEvent(buf *bytes.Buffer, m *Event) (err error) {
	return e.encodeMessageWithLeader(buf, &m.Message, 'E')
}

func (e *BinaryEncoder) EncodeHeartbeat(buf *bytes.Buffer, m *Heartbeat) (err error) {
	return e.encodeMessageWithLeader(buf, &m.Message, 'H')
}

func (e *BinaryEncoder) EncodeMetric(buf *bytes.Buffer, m *Metric) (err error) {
	return e.encodeMessageWithLeader(buf, &m.Message, 'M')
}
