package cat

import (
	"net"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/andywu1998/cat-go/message"
	"github.com/andywu1998/cat-go/test"
)

type brokenConn struct {
	net.Conn
	a *assert.Assertions
}

func (c *brokenConn) Write(b []byte) (n int, err error) {
	c.a.True(false, "This is a broken connection.")
	return
}

func (c *brokenConn) SetWriteDeadline(t time.Time) error {
	return errors.New("connection is broken")
}

var originConnInvalidHook = sender.messageDiscardHook

func receiveMessage(ch chan message.Messager, timeout time.Duration) message.Messager {
	select {
	case m := <-ch:
		return m
	case <-time.NewTimer(timeout).C:
		return nil
	}
}

func Test_Sender(t *testing.T) {
	defer func() {
		sender.messageDiscardHook = originConnInvalidHook
	}()

	a := assert.New(t)

	var ch = make(chan message.Messager, 10)

	sender.messageDiscardHook = func(m message.Messager) {
		ch <- m
	}

	go background(&sender)

	defer func() {
		scheduler.shutdownAndWaitGroup([]scheduleMixer{&sender})
	}()

	var m message.Messager

	// send a message without giving any active connection.
	// sender should be blocked at least 3 seconds.
	sender.high <- message.NewEvent("foo", "1", nil)
	m = receiveMessage(ch, senderBlockingTimeoutTime+time.Second)
	a.True(time.Now().Sub(m.GetTime()) > senderBlockingTimeoutTime)

	// should not block after the first 3s have elapsed
	sender.high <- message.NewEvent("foo", "2", nil)
	m = receiveMessage(ch, time.Second)
	a.True(time.Now().Sub(m.GetTime()) < time.Millisecond*10)

	// give a broken connection to sender.
	// the given connection should be discarded.
	sender.chConn <- &brokenConn{test.BlackHoleConn, a}

	// this message will be discarded immediately
	sender.high <- message.NewEvent("foo", "3", nil)
	m = receiveMessage(ch, time.Second)
	a.True(time.Now().Sub(m.GetTime()) < time.Millisecond*10)

	// a signal should be sent to router to make a new connection.
	signal := <-router.signals
	a.Equal(signalResetConnection, signal)

	// and the following message will be blocked again.
	sender.high <- message.NewEvent("foo", "4", nil)
	m = receiveMessage(ch, senderBlockingTimeoutTime+time.Second)
	a.True(time.Now().Sub(m.GetTime()) > senderBlockingTimeoutTime)

	// give a black hole connection to sender.
	// message should be serialized and sent immediately.
	sender.chConn <- test.BlackHoleConn
	sender.high <- message.NewEvent("foo", "5", nil)
	m = receiveMessage(ch, time.Second)
	a.Nil(m)

	// though the connection is broken, but it recovered immediately.
	// the 2nd message should not be discarded. (the 1st will because it triggered drop connection)
	sender.chConn <- &brokenConn{test.BlackHoleConn, a}
	sender.normal <- message.NewEvent("foo", "6", nil)
	m = receiveMessage(ch, time.Second)
	a.True(time.Now().Sub(m.GetTime()) < time.Millisecond*10)

	<-router.signals

	sender.chConn <- test.BlackHoleConn
	sender.normal <- message.NewEvent("foo", "7", nil)
	m = receiveMessage(ch, time.Second)
	a.Nil(m)
}
