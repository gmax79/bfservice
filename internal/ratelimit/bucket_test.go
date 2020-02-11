package ratelimit

import (
	"testing"
	"time"
)

func TestCreateInvalidBusket(t *testing.T) {
	_, err := CreateBucket(0, time.Second)
	if err == nil {
		t.Fatal("Basket cant be empty")
	}
	_, err = CreateBucket(5, time.Second*0)
	if err == nil {
		t.Fatal("Basket cant be 0 rated")
	}
}

func TestBusketMaxSize(t *testing.T) {
	busket, err := CreateBucket(10, time.Second)
	if err != nil {
		t.Fatal(err)
	}
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
	busket, err := CreateBucket(10, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i <= 10; i++ {
		if !busket.Score() {
			t.Fatal("Busket not be full")
		}
		time.Sleep(time.Millisecond * 120)
	}
	if !busket.Score() {
		t.Fatal("Busket not be full, was drained?")
	}
	if busket.Score() {
		t.Fatal("Busket must be full, fill drained item")
	}
}

func TestBusketIdleTime(t *testing.T) {
	list, err := CreateBucket(10, time.Millisecond*50)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond * 100)
	if list.Idletime() < time.Millisecond*100 {
		t.Error("Lifetime can not be less 100 ms")
	}
}
