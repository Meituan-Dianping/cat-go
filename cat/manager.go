package cat

import (
	"fmt"
	"math"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/andywu1998/cat-go/message"
)

type catMessageManager struct {
	index           uint32
	offset          uint32
	hour            int
	messageIdPrefix string

	flush message.Flush
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

func atomicLoadFloat64(f *float64) float64 {
	unsafeAddr := (*uint64)(unsafe.Pointer(f))
	atomic.LoadUint64(unsafeAddr)
	return math.Float64frombits(*unsafeAddr)
}

func atomicStoreFloat64(f *float64, v float64) {
	unsafeAddr := (*uint64)(unsafe.Pointer(f))
	atomic.StoreUint64(unsafeAddr, math.Float64bits(v))
}

func init() {
	manager.flush = func(m message.Messager) {
		switch m := m.(type) {
		case *message.Transaction:
			if m.Status != SUCCESS {
				sender.handleTransaction(m)
			} else if manager.hitSample(atomicLoadFloat64(&router.sample)) {
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
}
