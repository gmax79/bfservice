package buckets

import (
	"testing"
	"time"
)

func TestCreateTimeList(t *testing.T) {
	_, err := CreateTimeList(1)
	if err == nil {
		t.Fatal("List cant be len less 2")
	}
}

func TestLenEmptyTimeList(t *testing.T) {
	list, err := CreateTimeList(2)
	if err != nil {
		t.Fatal(err)
	}
	if list.Diff() != 0 {
		t.Fatal("Empty list must be len 0")
	}
}

func TestTimeListMaxSize(t *testing.T) {
	list, err := CreateTimeList(5)
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i <= 4; i++ {
		if list.Push() {
			t.Fatal("List not be full")
		}
	}
	if !list.Push() {
		t.Fatal("List must be full")
	}
	time.Sleep(time.Second * 2)
}

func TestTimeListDiff(t *testing.T) {
	list, err := CreateTimeList(11)
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i < 11; i++ {
		if list.Push() {
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
	_ = list.Push()
	time.Sleep(time.Millisecond * 100)
	if list.Lifetime() <= time.Millisecond*100 {
		t.Error("Duration can not be less 100 ms")
	}
}
