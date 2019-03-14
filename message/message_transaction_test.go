package message

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_TransactionDuration(t *testing.T) {
	var m = NewTransaction("foo", "bar", nil)

	const (
		beforeMs = 50
		before = beforeMs * time.Millisecond
	)

	var start = time.Now().Add(-before)
	m.durationStart = start

	m.Complete()

	assert.True(t, m.GetDuration() >= before && m.GetDuration() < before + 10 * time.Millisecond)
}

func Test_TransactionSetDuration(t *testing.T) {
	var m *Transaction

	const (
		duration = 50 * time.Millisecond
	)

	m = NewTransaction("foo", "bar", nil)
	m.SetDuration(duration)
	assert.Equal(t, duration, m.GetDuration())

	m = NewTransaction("foo", "bar", nil)
	m.SetDuration(duration)
	m.Complete()
	assert.Equal(t, duration, m.GetDuration())
}
