package cat

import (
	"fmt"
	"sync/atomic"
	"time"

	"../message"
)

type catMessageManager struct {
	index           uint32
	offset          uint32
	hour            int
	messageIdPrefix string
}

func (p *catMessageManager) sendTransaction(t *message.Transaction) {
	sender.handleTransaction(t)
}

func (p *catMessageManager) sendEvent(t *message.Event) {
	sender.handleEvent(t)
}

func (p *catMessageManager) flush(m message.Messager) {
	switch m := m.(type) {
	case *message.Transaction:
		if m.Status != SUCCESS {
			sender.handleTransaction(m)
		} else if p.hitSample(router.sample) {
			sender.handleTransaction(m)
		} else {
			aggregator.transaction.Put(m)
		}
	case *message.Event:
		if m.Status != SUCCESS {
			sender.handleEvent(m)
		} else {
			aggregator.event.Put(m)
		}
	default:
		logger.Warning("Unrecognized message type.")
	}
}

func (p *catMessageManager) hitSample(sampleRate float64) bool {
	if sampleRate > 1.0 {
		return true
	} else if sampleRate < 1e-9 {
		return false
	}
	var cycle = uint32(1 / sampleRate)

	var current, next uint32
	for {
		current = atomic.LoadUint32(&p.offset)
		next = (current + 1) % cycle
		if atomic.CompareAndSwapUint32(&p.offset, current, next) {
			break
		}
	}
	return next == 0
}

func (p *catMessageManager) nextId() string {
	hour := int(time.Now().Unix() / 3600)

	if hour != p.hour {
		p.hour = hour
		p.messageIdPrefix = fmt.Sprintf("%s-%s-%d", config.domain, config.ipHex, hour)

		currentIndex := atomic.LoadUint32(&p.index)
		if atomic.CompareAndSwapUint32(&p.index, currentIndex, 0) {
			logger.Info("MessageId prefix has changed to: %s", p.messageIdPrefix)
		}
	}

	return fmt.Sprintf("%s-%d", p.messageIdPrefix, atomic.AddUint32(&p.index, 1))
}

var manager = catMessageManager{
	index:  0,
	offset: 0,
	hour:   0,
}
