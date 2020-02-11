package ratelimit

import (
	"sync"
	"time"
)

// Limitation - object to limits elements check per size/minute
type Limitation struct {
	items        map[string]*Busket
	maxPerItem   int
	mutex        *sync.Mutex
	timeInterval time.Duration
}

// CreateLimitation - create limitation map - count per duration
func CreateLimitation(count int, duration time.Duration) *Limitation {
	var m Limitation
	m.items = make(map[string]*Busket)
	m.mutex = &sync.Mutex{}
	m.maxPerItem = count
	m.timeInterval = duration
	go func() {
		// garbage collector
		for {
			m.mutex.Lock()
			for k, t := range m.items {
				d := t.Idletime()
				if d > m.timeInterval {
					delete(m.items, k)
				}
			}
			m.mutex.Unlock()
			time.Sleep(time.Millisecond * 100)
		}
	}()
	return &m
}

// Check - check item for limitation
func (m *Limitation) Check(item string) (bool, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	v, ok := m.items[item]
	if !ok {
		var err error
		v, err = CreateBusket(m.maxPerItem, m.timeInterval)
		if err != nil {
			return false, err
		}
		m.items[item] = v
	}
	scored := v.Score()
	return scored, nil
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
