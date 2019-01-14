package cat

import (
	"sync"
)

type CCMap struct {
	count uint32
	buckets []CCMapBucket
	hasher CCMapHasher
}

type CCMapBucket struct {
	mu     sync.Mutex
	data   map[string]interface{}
}

type CCMapHasher func(name string) uint32

type CCMapCreator func(name string) interface{}

type CCMapComputer func(interface{}) error

func hasher(name string) uint32 {
	var h uint32 = 0
	if len(name) > 0 {
		for i:=0;i<len(name);i++ {
			h = 31 * h + uint32(name[i])
		}
	}
	return h
}

func NewCCMap(count int) *CCMap {
	var ccmap = &CCMap{
		count: uint32(count),
		buckets: make([]CCMapBucket, count),
		hasher: hasher,
	}
	for i := 0; i < count; i++ {
		ccmap.buckets[i] = CCMapBucket{
			mu:     sync.Mutex{},
			data:   make(map[string]interface{}),
		}
	}
	return ccmap
}

func (p *CCMap) compute(key string, creator CCMapCreator, computer CCMapComputer) {
	hash := p.hasher(key)
	slot := hash % p.count
	bucket := p.buckets[slot]

	bucket.mu.Lock()
	bucket.compute(key, creator, computer)
	defer bucket.mu.Unlock()
}

func (p *CCMapBucket) compute(key string, creator CCMapCreator, computer CCMapComputer) (err error) {
	var item interface{}
	var ok bool

	if item, ok = p.data[key]; !ok {
		p.data[key] = creator(key)
		item = p.data[key]
	}
	err = computer(item)
	return
}
