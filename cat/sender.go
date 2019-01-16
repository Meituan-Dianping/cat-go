package cat

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/Meituan-Dianping/cat-go/message"
)

var header = &message.Header{
	Domain:   config.domain,
	Hostname: config.hostname,
	Ip:       config.ip,

	MessageId:       "",
	ParentMessageId: "",
	RootMessageId:   "",
}

type catMessageSender struct {
	scheduleMixin

	normal  chan message.Messager
	high    chan message.Messager
	chConn  chan net.Conn
	encoder message.Encoder

	buf *bytes.Buffer

	conn net.Conn
}

func (s *catMessageSender) GetName() string {
	return "Sender"
}

func (s *catMessageSender) send(m message.Messager) {
	var buf = s.buf
	buf.Reset()

	if tree, ok := m.(*catMessageTree); ok {
		header.MessageId = tree.messageId
		header.MessageId = tree.parentMessageId
		header.MessageId = tree.rootMessageId
		m = &tree.Transaction
	} else {
		header.MessageId = manager.nextId()
		header.MessageId = ""
		header.MessageId = ""
	}

	if err := s.encoder.EncodeHeader(buf, header); err != nil {
		return
	}
	if err := s.encoder.EncodeMessage(buf, m); err != nil {
		return
	}

	fmt.Println("send")

	var b = make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(buf.Len()))

	if err := s.conn.SetWriteDeadline(time.Now().Add(time.Second * 3)); err != nil {
		logger.Warning("Error occurred while setting write deadline, connection has been dropped.")
		s.conn = nil
		router.signals <- signalResetConnection
	}

	if _, err := s.conn.Write(b); err != nil {
		logger.Warning("Error occurred while writing data, connection has been dropped.")
		s.conn = nil
		router.signals <- signalResetConnection
		return
	}
	if _, err := s.conn.Write(buf.Bytes()); err != nil {
		logger.Warning("Error occurred while writing data, connection has been dropped.")
		s.conn = nil
		router.signals <- signalResetConnection
		return
	}
	return
}

func (s *catMessageSender) handleMessageTree(tree *catMessageTree, hasProblem bool) {
	if hasProblem {
		select {
		case s.high <- tree:
		default:
			logger.Warning("High priority channel is full, transaction has been discarded.")
		}
	} else {
		select {
		case s.normal <- tree:
		default:
			logger.Warning("Normal priority channel is full, transaction has been discarded.")
		}
	}
}

func (s *catMessageSender) handleTransaction(trans *message.Transaction) {
	if trans.GetStatus() != SUCCESS {
		select {
		case s.high <- trans:
		default:
			logger.Warning("High priority channel is full, transaction has been discarded.")
		}
	} else {
		select {
		case s.normal <- trans:
		default:
			logger.Warning("Normal priority channel is full, transaction has been discarded.")
		}
	}
}

func (s *catMessageSender) handleEvent(event *message.Event) {
	select {
	case s.normal <- event:
	default:
		// logger.Warning("Normal priority channel is full, event has been discarded.")
	}
}

func (s *catMessageSender) beforeStop() {
	close(s.chConn)
	close(s.high)
	close(s.normal)

	for m := range s.high {
		s.send(m)
	}
	for m := range s.normal {
		s.send(m)
	}
}

func (s *catMessageSender) process() {
	if s.conn == nil {
		s.conn = <- s.chConn
		logger.Info("Received a new connection: %s", s.conn.RemoteAddr().String())
		return
	}

	select {
	case sig := <-s.signals:
		s.handle(sig)
	case conn := <-s.chConn:
		logger.Info("Received a new connection: %s", conn.RemoteAddr().String())
		s.conn = conn
	case m := <-s.high:
		// logger.Debug("Receive a message [%s|%s] from high priority channel", m.GetType(), m.GetName())
		s.send(m)
	case m := <-s.normal:
		// logger.Debug("Receive a message [%s|%s] from normal priority channel", m.GetType(), m.GetName())
		s.send(m)
	}
}

var sender = catMessageSender{
	scheduleMixin: makeScheduleMixedIn(signalSenderExit),
	normal:        make(chan message.Messager, normalPriorityQueueSize),
	high:          make(chan message.Messager, highPriorityQueueSize),
	chConn:        make(chan net.Conn),
	encoder:       message.NewBinaryEncoder(),
	buf:           bytes.NewBuffer([]byte{}),
	conn:          nil,
}
