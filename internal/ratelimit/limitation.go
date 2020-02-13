package ratelimit

import (
	"sync"
	"time"
)

const bucketsLifeTime = time.Second

// Limitation - object to limits elements check per size/minute
type Limitation struct {
	items      map[string]*Bucket
	maxPerItem int
	mutex      *sync.Mutex
	bf         func() *Bucket
}

// CreateLimitation - create limitation map - count per duration, requires function to create new buckets
func CreateLimitation(bf func() *Bucket) *Limitation {
	var m Limitation
	m.items = make(map[string]*Bucket)
	m.mutex = &sync.Mutex{}
	m.bf = bf
	clock := bf().clock
	go func() {
		// garbage collector
		for {
			m.mutex.Lock()
			for k, t := range m.items {
				d := t.Idletime()
				if d >= bucketsLifeTime {
					delete(m.items, k)
				}
			}
			m.mutex.Unlock()
			clock.Sleep(time.Millisecond * 100)
		}
	}()
	return &m
}

// Check - check item for limitation
func (m *Limitation) Check(item string) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	v, ok := m.items[item]
	if !ok {
		v = m.bf()
		m.items[item] = v
	}
	scored := v.Score()
	return scored
}

// Reset - remove item from limitation
func (m *Limitation) Reset(item string) bool {
	m.mutex.Lock()
	_, ok := m.items[item]
	delete(m.items, item)
	m.mutex.Unlock()
	return ok
}

// Size - return count of actual elements
func (m *Limitation) Size() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return len(m.items)
}
