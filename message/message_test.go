package message

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Complete(t *testing.T) {
	var m Messager

	m = NewEvent("foo", "bar", nil)
	m.Complete()

	m = NewHeartbeat("foo", "bar", nil)
	m.Complete()

	m = NewTransaction("foo", "bar", nil)
	m.Complete()
}

func Test_MessageTime(t *testing.T) {
	var m = NewMessage("foo", "bar", nil)

	nano := time.Now().UnixNano()
	nano -= nano % time.Millisecond.Nanoseconds()
	now := time.Unix(0, nano)

	m.SetTime(now)
	assert.Equal(t, now, m.GetTime())
}

func Test_EventConcurrentComplete(t *testing.T) {
	var messages = make([]Messager, 0)

	var flush = func(m Messager) {
		messages = append(messages, m)
	}

	var m = NewEvent("foo", "bar", flush)

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Complete()
		}()
	}
	wg.Wait()

	assert.Equal(t, 1, len(messages))
}
