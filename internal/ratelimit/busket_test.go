package ratelimit

import (
	"testing"
	"time"
)

func TestCreateInvalidBusket(t *testing.T) {
	_, err := CreateBusket(0, time.Second)
	if err == nil {
		t.Fatal("Basket cant be empty")
	}
	_, err = CreateBusket(5, time.Second*0)
	if err == nil {
		t.Fatal("Basket cant be 0 rated")
	}
}

func TestTimeListMaxSize(t *testing.T) {
	list, err := CreateTimeList(5)
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i <= 5; i++ {
		if list.Score() {
			t.Fatal("List not be full")
		}
	}
	if !list.Score() {
		t.Fatal("List must be full")
	}
	time.Sleep(time.Second * 2)
}

func TestTimeListDiff(t *testing.T) {
	list, err := CreateTimeList(10)
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i <= 10; i++ {
		if list.Score() {
			t.Fatal("List not be full")
		}
		time.Sleep(time.Millisecond * 120)
	}
	if list.Diff() <= time.Second {
		t.Error("Duration can not be less 1 second")
	}
}

func TestTimeListLifeTime(t *testing.T) {
	list, err := CreateTimeList(11)
	if err != nil {
		t.Fatal(err)
	}
	_ = list.Score()
	time.Sleep(time.Millisecond * 100)
	if list.Lifetime() <= time.Millisecond*100 {
		t.Error("Lifetime can not be less 100 ms")
	}
}
