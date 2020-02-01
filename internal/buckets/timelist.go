package buckets

import (
	"errors"
	"time"
)

// TimeList - limited list with time value
type TimeList struct {
	list []time.Time
	head int
	tail int
}

// CreateTimeList - create list for calcualate limits per duration
func CreateTimeList(size int) (*TimeList, error) {
	if size < 2 {
		return nil, errors.New("Timelist len can't be less 2")
	}
	var tl TimeList
	tl.list = make([]time.Time, size)
	tl.head = 0
	tl.tail = 0
	return &tl, nil
}

// Push - add now time in list, remove oldest if need
func (tl *TimeList) Push() bool {
	h := tl.head
	t := tl.tail
	now := time.Now()
	incmax := func(val int) int {
		if val = val + 1; val == len(tl.list) {
			val = 0
		}
		return val
	}
	if h == t {
		tl.list[h] = now
	}
	h = incmax(h)
	isfull := h == t
	if isfull {
		tl.tail = incmax(t)
	}
	tl.list[h] = now
	tl.head = h
	return isfull
}

// Diff - calculate time duration between head and tail
func (tl *TimeList) Diff() time.Duration {
	h := tl.head
	t := tl.tail
	if h == t {
		return 0
	}
	return tl.list[h].Sub(tl.list[t])
}

// Lifetime - calculate duration between now and list head
func (tl *TimeList) Lifetime() time.Duration {
	h := tl.head
	now := time.Now()
	return now.Sub(tl.list[h])
}
