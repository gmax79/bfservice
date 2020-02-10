package ratelimit

import (
	"time"
)

type Busket struct {
	capacity    float64
	fill        float64
	lastDrained time.Time
	ratems      int64
}

// CreateBusket - create object for limit per duration
func CreateBusket(size int, rate time.Duration) (*Busket, error) {
	var b Busket
	b.capacity = float64(size)
	b.fill = 0
	b.lastDrained = time.Now()
	b.ratems = rate.Milliseconds()
	return &b, nil
}

// Score - count event, add now time in list, remove oldest if need
func (b *Busket) Score() bool {
	b.drain()

	return false
}

// drain - remove fill level from busket per time
func (b *Busket) drain() {
	timeout := time.Now() - b.lastDrained
	timeout.
}

/*
// Diff - calculate time duration between head and tail
func (b *Busket) Diff() time.Duration {
	return 0
}

// Lifetime - calculate duration between now and list head
func (b *Busket) Lifetime() time.Duration {
	return 0
}
*/
