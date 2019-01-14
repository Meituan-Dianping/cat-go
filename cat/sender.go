package cat

import (
	"../message"
	"bytes"
	"encoding/binary"
	"net"
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
	signalsMixin

	normal  chan message.Messager
	high    chan message.Messager
	chConn  chan net.Conn
	encoder message.Encoder

	buf *bytes.Buffer

	conn net.Conn
}

func (sender *catMessageSender) GetName() string {
	return "Sender"
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
	if trans.GetStatus() != SUCCESS {
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
	for sender.isAlive {
		if sender.conn == nil {
			sender.conn = <-sender.chConn
			logger.Info("Received a new connection: %s", sender.conn.RemoteAddr().String())
		} else {
			select {
			case signal := <-sender.signals:
				if signal == signalShutdown {
					close(sender.chConn)
					close(sender.high)
					close(sender.normal)
					sender.stop()
				}
			case conn := <-sender.chConn:
				logger.Info("Received a new connection: %s", conn.RemoteAddr().String())
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

	for m := range sender.high {
		sender.send(m)
	}
	for m := range sender.normal {
		sender.send(m)
	}
	sender.exit()
}

var sender = catMessageSender{
	signalsMixin: makeSignalsMixedIn(signalSenderExit),
	normal:       make(chan message.Messager, normalPriorityQueueSize),
	high:         make(chan message.Messager, highPriorityQueueSize),
	chConn:       make(chan net.Conn),
	encoder:      message.NewBinaryEncoder(),
	buf:          bytes.NewBuffer([]byte{}),
	conn:         nil,
}
