package cat

import (
	"fmt"
	"sync/atomic"
	"time"

	"../message"
)

type catMessageManager struct {
	index  uint32
	offset uint32
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
		} else if p.isSample() {
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

func (p *catMessageManager) isSample() bool {
	if router.sample > 1.0 {
		return true
	} else if router.sample < 1e-9 {
		return false
	}
	var cycle = uint32(1 / router.sample)

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
	// TODO reset every hour.
	hour := time.Now().Unix() / 3600
	return fmt.Sprintf("%s-%s-%d-%d", config.domain, config.ipHex, hour, atomic.AddUint32(&p.index, 1))
}

var manager = catMessageManager{
	index:  0,
	offset: 0,
}
