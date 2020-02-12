package ratelimit

import (
	"testing"
	"time"

	"github.com/jdeal-mediamath/clockwork"
)

func TestCreateInvalidBusket(t *testing.T) {
	clock := clockwork.NewFakeClock()
	_, err := CreateBucketsFactory(0, time.Second, clock)
	if err == nil {
		t.Fatal("Basket cant be empty")
	}
	_, err = CreateBucketsFactory(5, time.Second*0, clock)
	if err == nil {
		t.Fatal("Basket cant be zero rated")
	}
}

func TestBusketMaxSize(t *testing.T) {
	clock := clockwork.NewFakeClock()
	bf, err := CreateBucketsFactory(10, time.Second, clock)
	if err != nil {
		t.Fatal(err)
	}
	busket := bf()
	for i := 1; i <= 10; i++ {
		if !busket.Score() {
			t.Fatal("Busket not be full")
		}
	}
	if busket.Score() {
		t.Fatal("Busket must be full")
	}
}

func TestBusketRating(t *testing.T) {
	clock := clockwork.NewFakeClock()
	bf, err := CreateBucketsFactory(10, time.Second, clock)
	if err != nil {
		t.Fatal(err)
	}
	busket := bf()
	for i := 1; i <= 10; i++ {
		if !busket.Score() {
			t.Fatal("Busket not be full")
		}
		clock.Advance(time.Millisecond * 120)
	}
	if !busket.Score() {
		t.Fatal("Busket not be full, was drained?")
	}
	if busket.Score() {
		t.Fatal("Busket must be full, fill drained item")
	}
}

func TestBusketIdleTime(t *testing.T) {
	clock := clockwork.NewFakeClock()
	bf, err := CreateBucketsFactory(10, time.Millisecond*50, clock)
	if err != nil {
		t.Fatal(err)
	}
	busket := bf()
	clock.Advance(time.Millisecond * 100)
	if busket.Idletime() < time.Millisecond*100 {
		t.Error("Lifetime can not be less 100 ms")
	}
}
