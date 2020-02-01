package buckets

import (
	"sync"
	"time"
)

// Limitation - object to limits elements check per size/minute
type Limitation struct {
	items        map[string]*TimeList
	maxPerItem   int
	mutex        *sync.Mutex
	timeInterval time.Duration
}

// CreateLimitation - create limitation map - count per duration
func CreateLimitation(count int, duration time.Duration) *Limitation {
	var m Limitation
	m.items = make(map[string]*TimeList)
	m.mutex = &sync.Mutex{}
	m.maxPerItem = count
	m.timeInterval = duration
	go func() {
		// garbage collector
		m.mutex.Lock()
		for k, t := range m.items {
			if t.Lifetime() > m.timeInterval {
				delete(m.items, k)
			}
		}
		m.mutex.Unlock()
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
		m.items[item], err = CreateTimeList(m.maxPerItem)
		if err != nil {
			return false, err
		}
	}
	if v.Push() && v.Diff() > m.timeInterval {
		return false, nil
	}
	return true, nil
}
