package test

import (
	"net"
	"time"
)

type blackHoleConn struct {
}

var BlackHoleConn = &blackHoleConn{}

func (c *blackHoleConn) Read(b []byte) (n int, err error) {
	return
}

func (c *blackHoleConn) Write(b []byte) (n int, err error) {
	return
}

func (c *blackHoleConn) Close() error {
	return nil
}

func (c *blackHoleConn) LocalAddr() net.Addr {
	return &net.IPAddr{}
}

func (c *blackHoleConn) RemoteAddr() net.Addr {
	return &net.IPAddr{}
}

func (c *blackHoleConn) SetDeadline(t time.Time) error {
	return nil
}

func (c *blackHoleConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *blackHoleConn) SetWriteDeadline(t time.Time) error {
	return nil
}
