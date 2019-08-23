package message

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

const (
	txtTAB = "\t"
	txtLF = "\n"
)

type TxtEncoder struct {
	encoderBase
}

func NewTxtEncoder() *TxtEncoder {
	return &TxtEncoder{}
}

func (e *TxtEncoder) writeString(buf *bytes.Buffer, s string) (err error) {

	if _, err = buf.WriteString(s+txtTAB); err != nil {
		return
	}
	return
}

func (e *TxtEncoder) EncodeHeader(buf *bytes.Buffer, header *Header) (err error) {
	if err = e.writeString(buf, ReadableProtocol); err != nil {
		return
	}

	//buf.WriteString(txtTAB)
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
	buf.WriteString("")

	//追加换行
	buf.WriteString(txtLF)
	return
}

func (e *TxtEncoder) encodeMessageStart(buf *bytes.Buffer, m Messager) (err error) {
	var timestampTemp = m.GetTime().UnixNano() / time.Millisecond.Nanoseconds()

	micTime := timestampTemp % 1000
	formatTime := m.GetTime().Format("2006-01-02 15:04:05")
	var timestamp = fmt.Sprintf("%s.%03d", formatTime, micTime)


	if err = e.writeString(buf, timestamp); err != nil {
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

func (e *TxtEncoder) encodeMessageEnd(buf *bytes.Buffer, m Messager) (err error) {
	if err = e.writeString(buf, m.GetStatus()); err != nil {
		return
	}

	if m.GetData() == nil {
		if err = e.writeString(buf, ""); err != nil {
			return
		}
	} else {

		if err = e.writeString(buf, m.GetData().String()); err != nil {
			return
		}
	}
	return
}

func (e *TxtEncoder) encodeMessageWithLeader(buf *bytes.Buffer, m *Message, leader rune) (err error) {
	if _, err = buf.WriteRune(leader); err != nil {
		return
	}
	if err = e.encodeMessageStart(buf, m); err != nil {
		return
	}
	if err = e.encodeMessageEnd(buf, m); err != nil {
		return
	}
	//追加换行
	buf.WriteString(txtLF)

	return
}

func (e *TxtEncoder) EncodeMessage(buf *bytes.Buffer, message Messager) (err error) {
	return encodeMessage(e, buf, message)
}

func (e *TxtEncoder) EncodeTransaction(buf *bytes.Buffer, trans *Transaction) (err error) {
	if _, err = buf.WriteRune('t'); err != nil {
		return
	}
	if err = e.encodeMessageStart(buf, trans); err != nil {
		return
	}

	//追加换行
	buf.WriteString(txtLF)

	//子节点
	for _, message := range trans.GetChildren() {
		if err = e.EncodeMessage(buf, message); err != nil {
			return
		}
	}

	if _, err = buf.WriteRune('T'); err != nil {
		return
	}

	//结束时间
	timestampTemp := trans.GetTime().UnixNano() / time.Millisecond.Nanoseconds()
	micTime := timestampTemp % 1000
	formatTime := trans.GetTime().Format("2006-01-02 15:04:05")
	timestamp := fmt.Sprintf("%s.%03d", formatTime, micTime)

	if err = e.writeString(buf, timestamp); err != nil {
		return
	}
	if err = e.writeString(buf, trans.GetType()); err != nil {
		return
	}
	if err = e.writeString(buf, trans.GetName()); err != nil {
		return
	}

	//状态
	if err = e.writeString(buf, trans.GetStatus()); err != nil {
		return
	}

	//持续时间
	duration := trans.GetDuration().Nanoseconds() / time.Microsecond.Nanoseconds()
	if err = e.writeString(buf, strconv.FormatInt(duration, 10)+"us"); err != nil {
		return
	}

	//数据内容
	if trans.GetData() == nil {
		if err = e.writeString(buf, ""); err != nil {
			return
		}
	} else {

		if err = e.writeString(buf, trans.GetData().String()); err != nil {
			return
		}
	}

	//追加换行
	buf.WriteString(txtLF)

	return
}

func (e *TxtEncoder) EncodeEvent(buf *bytes.Buffer, m *Event) (err error) {
	return e.encodeMessageWithLeader(buf, &m.Message, 'E')
}

func (e *TxtEncoder) EncodeHeartbeat(buf *bytes.Buffer, m *Heartbeat) (err error) {
	//return err
	return e.encodeMessageWithLeader(buf, &m.Message, 'H')
}

func (e *TxtEncoder) EncodeMetric(buf *bytes.Buffer, m *Metric) (err error) {
	return e.encodeMessageWithLeader(buf, &m.Message, 'M')
}
