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
		if l.Check(testid) == false {
			t.Fatal("Limitation must non-blocking")
		}
	}
	for i := 1; i <= 10; i++ {
		if l.Check(testid) == true {
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
		if l.Check(testid) == false {
			t.Fatal("Limitation must non-blocking")
		}
	}
	// wait, free element by gc
	clock.Advance(time.Millisecond * 100)
	if l.Check(testid) == false {
		t.Fatal("Limitation must non-blocking")
	}
	if l.Check(testid) == true {
		t.Fatal("Limitation must blocking")
	}

	// wait time to drain all elements
	clock.Advance(time.Second)

	// sleep to run gc goroutine
	time.Sleep(time.Millisecond * 100)

	if l.Size() != 0 {
		t.Fatal("Limitation must be garbaged")
	}
}

func TestLimitationReset(t *testing.T) {
	clock := clockwork.NewFakeClock()
	bf, err := CreateBucketsFactory(10, time.Second, clock)
	if err != nil {
		t.Fatal(err)
	}
	l := CreateLimitation(bf)
	l.Check(testid)
	found := l.Reset(testid)
	if !found {
		t.Fatal("bucket not found ?")
	}
}
