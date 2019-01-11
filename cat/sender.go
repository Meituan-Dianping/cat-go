package cat

import (
	"../message"
	"bytes"
	"encoding/binary"
	"net"
)

type catMessageSender struct {
	normal  chan message.Messager
	high    chan message.Messager
	chConn  chan net.Conn
	encoder message.Encoder

	buf *bytes.Buffer

	conn net.Conn
}

func createHeader() *message.Header {
	return &message.Header{
		Domain:   config.domain,
		Hostname: config.hostname,
		Ip:       config.ip,

		MessageId:       manager.nextId(),
		ParentMessageId: "",
		RootMessageId:   "",
	}
}

func (sender *catMessageSender) send(m message.Messager) (err error) {
	var buf = sender.buf
	buf.Reset()

	var header = createHeader()
	sender.encoder.EncodeHeader(buf, header)
	sender.encoder.EncodeMessage(buf, m)

	var b = make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(buf.Len()))

	if _, err = sender.conn.Write(b); err != nil {
		return
	}
	if _, err = sender.conn.Write(buf.Bytes()); err != nil {
		return
	}
	return
}

func (sender *catMessageSender) handleTransaction(trans *message.Transaction) {
	if trans.GetStatus() != CAT_SUCCESS {
		select {
		case sender.high <- trans:
		default:
			logger.Warning("High priority channel is full, transaction has been discarded.")
		}
	} else {
		select {
		case sender.normal <- trans:
		default:
			logger.Warning("Normal priority channel is full, transaction has been discarded.")
		}
	}
}

func (sender *catMessageSender) handleEvent(event *message.Event) {
	select {
	case sender.normal <- event:
	default:
		logger.Warning("Normal priority channel is full, event has been discarded.")
	}
}

func (sender *catMessageSender) Background() {
	for {
		if sender.conn == nil {
			sender.conn = <-sender.chConn
			logger.Info("Received a new connection: %s", sender.conn.LocalAddr().String())
		} else {
			select {
			case conn := <-sender.chConn:
				logger.Info("Received a new connection: %s", conn.LocalAddr().String())
				sender.conn = conn
			case m := <-sender.high:
				// logger.Debug("Receive a message [%s|%s] from high priority channel", m.GetType(), m.GetName())
				sender.send(m)
			case m := <-sender.normal:
				// logger.Debug("Receive a message [%s|%s] from normal priority channel", m.GetType(), m.GetName())
				sender.send(m)
			}
		}
	}
}

var sender = catMessageSender{
	normal:  make(chan message.Messager, NORMAL_PRIORITY_QUEUE_SIZE),
	high:    make(chan message.Messager, HIGH_PRIORITY_QUEUE_SIZE),
	chConn:  make(chan net.Conn),
	encoder: message.NewBinaryEncoder(),
	buf:     bytes.NewBuffer([]byte{}),
	conn:    nil,
}
