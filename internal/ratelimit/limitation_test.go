package ratelimit

import (
	"testing"
	"time"
)

const testid = "dummy"

func TestLimitationInvalidCreation(t *testing.T) {
	l := CreateLimitation(1, time.Second)
	_, err := l.Check(testid)
	if err == nil {
		t.Fatal("Must be error, cant create limitation with len < 2")
	}
}

func TestLimitationFast(t *testing.T) {
	l := CreateLimitation(10, time.Second)
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
	l := CreateLimitation(10, time.Second)
	for i := 1; i <= 10; i++ {
		state, err := l.Check(testid)
		if err != nil {
			t.Fatal(err)
		}
		if state == false {
			t.Fatal("Limitation must non-blocking")
		}
		time.Sleep(time.Millisecond * 100)
	}
	// wait, free some elements by gc
	time.Sleep(time.Millisecond * 100)
	state, err := l.Check(testid)
	if err != nil {
		t.Fatal(err)
	}
	if state == false {
		t.Fatal("Limitation must non-blocking")
	}
	time.Sleep(time.Millisecond * 100)
	state, err = l.Check(testid)
	if err != nil {
		t.Fatal(err)
	}
	if state == false {
		t.Fatal("Limitation must non-blocking")
	}
	// wait lifetime,
	time.Sleep(time.Millisecond * 1100)
	if l.Size() != 0 {
		t.Fatal("timelist must be garbaged")
	}
}
