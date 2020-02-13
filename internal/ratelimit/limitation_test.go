package ratelimit

import (
	"testing"
	"time"

	"github.com/jdeal-mediamath/clockwork"
)

const testid = "dummy"

func TestLimitationFast(t *testing.T) {
	clock := clockwork.NewFakeClock()
	bf, err := CreateBucketsFactory(10, time.Second, clock)
	if err != nil {
		t.Fatal(err)
	}
	l := CreateLimitation(bf)
	for i := 1; i <= 10; i++ {
		state, err := l.Check(testid)
		if err != nil {
			t.Fatal(err)
		}
		if state == false {
			t.Fatal("Limitation must non-blocking")
		}
	}
	for i := 1; i <= 10; i++ {
		state, err := l.Check(testid)
		if err != nil {
			t.Fatal(err)
		}
		if state == true {
			t.Fatal("Limitation must be blocked")
		}
	}
}

func TestLimitationGC(t *testing.T) {
	clock := clockwork.NewFakeClock()
	bf, err := CreateBucketsFactory(10, time.Second, clock)
	if err != nil {
		t.Fatal(err)
	}
	l := CreateLimitation(bf)
	for i := 1; i <= 10; i++ {
		state, err := l.Check(testid)
		if err != nil {
			t.Fatal(err)
		}
		if state == false {
			t.Fatal("Limitation must non-blocking")
		}
	}
	// wait, free element by gc
	clock.Advance(time.Millisecond * 100)
	state, err := l.Check(testid)
	if err != nil {
		t.Fatal(err)
	}
	if state == false {
		t.Fatal("Limitation must non-blocking")
	}
	state, err = l.Check(testid)
	if err != nil {
		t.Fatal(err)
	}
	if state == true {
		t.Fatal("Limitation must blocking")
	}

	// wait lifetime of bucket
	clock.Advance(bucketsLifeTime)

	// sleep, run gc goroutine
	time.Sleep(time.Millisecond * 100)

	if l.Size() != 0 {
		t.Fatal("Limitation must be garbaged")
	}
}
