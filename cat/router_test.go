package cat

import (
	"testing"
	"time"

	"github.com/andywu1998/cat-go/test"
)

func Test_Router(t *testing.T) {
	a := test.NewAssertion(t)

	manager := newTestServerManager()

	routerInterval = time.Second
	config.httpServerAddresses = []serverAddress{
		{"127.0.0.1", 2080},
	}

	go background(&router)
	time.Sleep(time.Second)

	defer func() {
		routerInterval = defaultRouterInterval
		config.httpServerAddresses = []serverAddress{}
		scheduler.shutdownAndWaitGroup([]scheduleMixer{&router})
		manager.shutdownAll()
	}()

	// test startup while router service is invalid.
	a.False(IsEnabled())

	// router service has started, but no valid remote service.
	manager.startHTTP(8080)
	time.Sleep(time.Second)
	a.False(IsEnabled())

	// remote server has started, but conn has not been received by sender.
	manager.startHTTP(2280)
	manager.startHTTP(2281)
	time.Sleep(time.Second)
	a.False(IsEnabled())

	// test get connection
	a.ReceiveConn(sender.chConn, "127.0.0.1:2280")
	a.True(IsEnabled())

	// test modify sampleRate
	manager.sampleRate = 0.5
	time.Sleep(time.Second)

	go func() {
		a.Equal(0.5, router.sample)
	}()
	time.Sleep(time.Millisecond * 10)

	// test toggle blocking
	a.True(IsEnabled())
	manager.block = true
	time.Sleep(time.Second)

	go func() {
		a.False(IsEnabled())
	}()
	time.Sleep(time.Millisecond * 10)

	manager.block = false
	time.Sleep(time.Second)

	go func() {
		a.True(IsEnabled())
	}()
	time.Sleep(time.Millisecond * 10)

	// test router change
	manager.ports = []int{2281, 2280, 2282}
	time.Sleep(time.Second)

	a.ReceiveConn(sender.chConn, "127.0.0.1:2281")
	time.Sleep(time.Millisecond * 10)

	// test router change, but first of 2 are invalid
	manager.ports = []int{2285, 2284, 2280}
	time.Sleep(time.Second)

	a.ReceiveConn(sender.chConn, "127.0.0.1:2280")
	time.Sleep(time.Millisecond * 10)

	// test router change, but all of them are invalid, then the first becomes valid.
	manager.ports = []int{2283, 2284, 2285}
	time.Sleep(time.Second)
	a.True(IsEnabled())

	// shouldn't receive any connection in this case.
	a.ReceiveConnN(sender.chConn)

	manager.startHTTP(2283)
	a.ReceiveConn(sender.chConn, "127.0.0.1:2283")
	time.Sleep(time.Millisecond * 10)

	// conn has been discarded by sender, a new connection should be established.
	sender.resetConnection()

	a.ReceiveConn(sender.chConn, "127.0.0.1:2283")
	time.Sleep(time.Millisecond * 10)
}
