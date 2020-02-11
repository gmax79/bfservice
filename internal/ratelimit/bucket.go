package ratelimit

import (
	"errors"
	"math"
	"time"
)

// Bucket - main object to implentation rate limit algorithm
type Bucket struct {
	capacity         int
	fill             int
	lastDrained      time.Time
	ratems           float64
	lastDrainResidue float64
}

// CreateBucket - create object for limit per duration
func CreateBucket(size int, rate time.Duration) (*Bucket, error) {
	if size <= 0 {
		return nil, errors.New("invalid size parameter")
	}
	if rate <= 0 {
		return nil, errors.New("invalid rate parameter")
	}
	var b Bucket
	b.capacity = size
	b.fill = 0
	b.lastDrained = time.Now()
	b.ratems = float64(rate.Milliseconds())
	b.lastDrainResidue = 0
	return &b, nil
}

// Score - count event, add now time in list, remove oldest if need
func (b *Bucket) Score() bool {
	b.drain()
	if b.fill == b.capacity {
		return false // busket empty
	}
	b.fill++    // count
	return true // pass
}

// drain - remove fill level from busket with configured rate
func (b *Bucket) drain() {
	now := time.Now()
	timeoutms := now.Sub(b.lastDrained).Milliseconds()
	drainedCount := b.lastDrainResidue + float64(timeoutms)/b.ratems
	if drainedCount >= 1 {
		b.lastDrained = now
		drainedCount, b.lastDrainResidue = math.Modf(drainedCount)
		todrain := int(drainedCount)
		if todrain > b.fill {
			todrain = b.fill
		}
		b.fill -= todrain
	}
}

// Idletime - calculate idle time for gc
func (b *Bucket) Idletime() time.Duration {
	now := time.Now()
	return now.Sub(b.lastDrained)
}
