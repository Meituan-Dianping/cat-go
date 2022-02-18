package cat

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"

	"github.com/andywu1998/cat-go/message"
)

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

type catMessageSender struct {
	scheduleMixin

	normal  chan message.Messager
	high    chan message.Messager
	chConn  chan net.Conn
	encoder message.Encoder

	buf *bytes.Buffer

	conn net.Conn

	lastActiveTime time.Time

	messageDiscardHook func(m message.Messager)
}

func (s *catMessageSender) GetName() string {
	return "Sender"
}

func (s *catMessageSender) resetConnection() {
	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			logger.Warning("Error occurred while closing current connection.")
		} else {
			logger.Info("Broken connection %s has been closed", s.conn.RemoteAddr().String())
		}
		s.conn = nil
	}
	router.signals <- signalResetConnection
}

func (s *catMessageSender) send(m message.Messager) {
	if s.conn == nil {
		s.messageDiscardHook(m)
		return
	}

	var buf = s.buf
	buf.Reset()

	var header = createHeader()
	if err := s.encoder.EncodeHeader(buf, header); err != nil {
		return
	}
	if err := s.encoder.EncodeMessage(buf, m); err != nil {
		return
	}

	var b = make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(buf.Len()))

	if err := s.conn.SetWriteDeadline(time.Now().Add(defaultWriteDeadline)); err != nil {
		logger.Warning("Error occurred while setting write deadline, connection has been dropped.")
		s.messageDiscardHook(m)
		s.resetConnection()
		return
	}

	if _, err := s.conn.Write(b); err != nil {
		logger.Warning("Error occurred while writing data, connection has been dropped. %s", err)
		s.messageDiscardHook(m)
		s.resetConnection()
		return
	}
	if _, err := s.conn.Write(buf.Bytes()); err != nil {
		logger.Warning("Error occurred while writing data, connection has been dropped. %s", err)
		s.messageDiscardHook(m)
		s.resetConnection()
		return
	}
	return
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
			// logger.Warning("Normal priority channel is full, transaction has been discarded.")
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

func (s *catMessageSender) afterStart() {
	s.lastActiveTime = time.Now()
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
		if time.Now().Sub(s.lastActiveTime) < senderBlockingTimeoutTime {
			logger.Info("Sender is blocking wait for a new connection")
			select {
			case s.conn = <-s.chConn:
				logger.Info("Received a connection: %s", s.conn.RemoteAddr().String())
			case <-time.NewTimer(senderBlockingTimeoutTime).C:
				logger.Warning("Can't get a new connection, further messages will be discarded.")
			}
			return
		}
	} else {
		s.lastActiveTime = time.Now()
	}

	select {
	case sig := <-s.signals:
		s.handle(sig)
	case conn := <-s.chConn:
		logger.Info("Connection switch to: %s", conn.RemoteAddr().String())
		if s.conn != nil {
			if err := s.conn.Close(); err != nil {
				logger.Warning("Error occurred while closing previous connection")
			} else {
				logger.Info("Previous connection %s has been closed", s.conn.RemoteAddr().String())
			}
		}
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
	scheduleMixin:   makeScheduleMixedIn(signalSenderExit),
	normal:          make(chan message.Messager, normalPriorityQueueSize),
	high:            make(chan message.Messager, highPriorityQueueSize),
	chConn:          make(chan net.Conn),
	encoder:         message.NewBinaryEncoder(),
	buf:             bytes.NewBuffer([]byte{}),
	conn:            nil,
	messageDiscardHook: func(messager message.Messager) {
		logger.Warning("Message has been discarded due to no active connection")
	},
}
