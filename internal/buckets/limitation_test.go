package buckets

import (
	"testing"
	"time"
)

const testid = "dummy"

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

}
