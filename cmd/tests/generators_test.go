package main

import (
	"testing"

	"github.com/gmax79/bfservice/internal/netsupport"
)

func TestRandomStringGenerator(t *testing.T) {
	gen := stringGenerator(10, 4, "check")
	counter := 0
	empty := 0
	for i := 0; i < 20; i++ {
		login := gen()
		if login == "check" {
			counter++
		}
		if login == "" {
			empty++
		}
	}
	if counter != 4 || empty != 6 {
		t.Fatal("invalid result by stringGenerator")
	}
}

func TestRandomIPGenerator(t *testing.T) {
	gen := ipGenerator(10, "192.168.1.0", 4, "192.168.1.1")
	counter := 0
	empty := 0
	for i := 0; i < 20; i++ {
		ip := gen()
		if ip == "192.168.1.1" {
			counter++
		}
		if ip == "" {
			empty++
		} else {
			var a netsupport.IPAddr
			err := a.Parse(ip)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
	if counter != 4 || empty != 6 {
		t.Fatal("invalid result by ipGenerator")
	}
}
