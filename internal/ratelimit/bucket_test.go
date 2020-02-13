package ratelimit

import (
	"testing"
	"time"

	"github.com/jdeal-mediamath/clockwork"
)

func TestCreateInvalidBucket(t *testing.T) {
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

func TestBucketMaxSize(t *testing.T) {
	clock := clockwork.NewFakeClock()
	bf, err := CreateBucketsFactory(10, time.Second, clock)
	if err != nil {
		t.Fatal(err)
	}
	bucket := bf()
	for i := 1; i <= 10; i++ {
		if !bucket.Score() {
			t.Fatal("Bucket not be full")
		}
	}
	if bucket.Score() {
		t.Fatal("Bucket must be full")
	}
}

func TestBucketRating(t *testing.T) {
	clock := clockwork.NewFakeClock()
	bf, err := CreateBucketsFactory(10, time.Second, clock)
	if err != nil {
		t.Fatal(err)
	}
	bucket := bf()
	for i := 1; i <= 10; i++ {
		if !bucket.Score() {
			t.Fatal("Bucket not be full")
		}
		clock.Advance(time.Millisecond * 100) // Drain element
	}
	if bucket.GetScores() != 0 {
		t.Fatal("bucket must be empty, all elements drained")
	}
}

func TestBucketRating2(t *testing.T) {
	clock := clockwork.NewFakeClock()
	bf, err := CreateBucketsFactory(10, time.Second, clock)
	if err != nil {
		t.Fatal(err)
	}
	bucket := bf()
	for i := 1; i <= 10; i++ {
		if !bucket.Score() {
			t.Fatal("Bucket not be full")
		}
	}
	clock.Advance(time.Millisecond * 500) // Drain 5 elements (0.5 sec * 10elem/sec)
	if bucket.GetScores() != 5 {
		t.Fatal("bucket must be 5 elements")
	}
}

func TestBucketIdleTime(t *testing.T) {
	clock := clockwork.NewFakeClock()
	bf, err := CreateBucketsFactory(10, time.Millisecond*50, clock)
	if err != nil {
		t.Fatal(err)
	}
	bucket := bf()
	clock.Advance(time.Millisecond * 100)
	if bucket.Idletime() < time.Millisecond*100 {
		t.Error("Lifetime can not be less 100 ms")
	}
}

func TestBucketFullDraining(t *testing.T) {
	clock := clockwork.NewFakeClock()
	bf, err := CreateBucketsFactory(10, time.Second, clock)
	if err != nil {
		t.Fatal(err)
	}
	bucket := bf()
	for i := 1; i <= 10; i++ {
		if !bucket.Score() {
			t.Fatal("Bucket not be full")
		}
	}
	clock.Advance(time.Second * 2)
	if bucket.GetScores() != 0 {
		t.Fatal("Bucket must be empty")
	}

}
