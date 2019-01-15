package cat

import (
	"sync"
)

type ccMap struct {
	count   uint32
	buckets []ccMapBucket
	hasher  ccMapHasher
}

type ccMapBucket struct {
	mu   sync.Mutex
	data map[string]interface{}
}

type ccMapHasher func(name string) uint32

type ccMapCreator func(name string) interface{}

type ccMapComputer func(interface{}) error

func hasher(name string) uint32 {
	var h uint32 = 0
	if len(name) > 0 {
		for i := 0; i < len(name); i++ {
			h = 31*h + uint32(name[i])
		}
	}
	return h
}

func newCCMap(count int) *ccMap {
	var ccmap = &ccMap{
		count:   uint32(count),
		buckets: make([]ccMapBucket, count),
		hasher:  hasher,
	}
	for i := 0; i < count; i++ {
		ccmap.buckets[i] = ccMapBucket{
			mu:   sync.Mutex{},
			data: make(map[string]interface{}),
		}
	}
	return ccmap
}

func (p *ccMap) compute(key string, creator ccMapCreator, computer ccMapComputer) (err error) {
	hash := p.hasher(key)
	slot := hash % p.count
	bucket := p.buckets[slot]

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	return bucket.compute(key, creator, computer)
}

func (p *ccMapBucket) compute(key string, creator ccMapCreator, computer ccMapComputer) (err error) {
	var item interface{}
	var ok bool

	if item, ok = p.data[key]; !ok {
		p.data[key] = creator(key)
		item = p.data[key]
	}
	err = computer(item)
	return
}
