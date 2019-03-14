package test

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

type assertions struct {
	assert.Assertions
}

func NewAssertion(t *testing.T) *assertions {
	return &assertions{*assert.New(t)}
}

func (a *assertions) ReceiveConn(ch chan net.Conn, expectAddr string) {
	select {
	case conn := <-ch:
		a.Equal(expectAddr, conn.RemoteAddr().String())
		if err := conn.Close(); err != nil {
			panic(err)
		}
	case <-time.NewTimer(time.Second).C:
		a.True(false, "Can't receive a connection within 1 second")
	}
}

func (a *assertions) ReceiveConnN(ch chan net.Conn) {
	select {
	case conn := <-ch:
		a.True(false, "Should not receive a connection in this case")
		if err := conn.Close(); err != nil {
			panic(err)
		}
	case <-time.NewTimer(time.Second).C:
		a.True(true)
	}
}
