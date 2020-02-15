package ratelimit

import (
	"errors"
	"math"
	"time"

	"github.com/jdeal-mediamath/clockwork"
)

// Bucket - main object to implementation rate limit algorithm
type Bucket struct {
	capacity         int
	fill             int
	lastDrained      time.Time
	ratems           float64
	lastDrainResidue float64
	clock            clockwork.Clock
}

// CreateBucketsFactory - initialize function-factory to instanitiate new buckets - size per rate
func CreateBucketsFactory(size int, rate time.Duration, clock clockwork.Clock) (func() *Bucket, error) {
	if size <= 0 {
		return nil, errors.New("invalid size parameter")
	}
	if rate <= 0 {
		return nil, errors.New("invalid rate parameter")
	}
	return func() *Bucket {
		var b Bucket
		b.clock = clock
		b.capacity = size
		b.fill = 0
		b.lastDrained = b.clock.Now()
		b.ratems = float64(size) / float64(rate.Milliseconds())
		b.lastDrainResidue = 0
		return &b
	}, nil
}

// Score - count event, add now time in list, remove oldest if need
func (b *Bucket) Score() bool {
	b.drain()
	if b.fill == b.capacity {
		return false // bucket full, mpt scored
	}
	b.fill++    // scored
	return true // pass
}

// GetScores - return fill level of bucket in now time
func (b *Bucket) GetScores() int {
	b.drain()
	return b.fill
}

// Empty - check bucket empty state
func (b *Bucket) Empty() bool {
	return b.GetScores() == 0
}

// Drain - remove fill level from busket with configured rate
func (b *Bucket) drain() {
	now := b.clock.Now()
	timeoutms := now.Sub(b.lastDrained).Milliseconds()
	drainedCount := b.lastDrainResidue + float64(timeoutms)*b.ratems
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
	now := b.clock.Now()
	return now.Sub(b.lastDrained)
}
