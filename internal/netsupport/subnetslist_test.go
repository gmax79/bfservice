package netsupport

import (
	"log"
	"strconv"
	"testing"

	"github.com/gmax79/bfservice/internal/storage"
)

func makesubnet(s string) Subnet {
	var snet Subnet
	err := snet.Parse(s)
	if err != nil {
		log.Fatal(err)
	}
	return snet
}

func createSet() (storage.SetProvider, error) {
	memstor, err := storage.ConnectMemory()
	if err != nil {
		return err
	}
	set, err := memstor.CreateSet("test")
	if err != nil {
		return err
	}
	return set, nil
}

func TestSubnetsListInList(t *testing.T) {
	set, err := createSet()
	if err != nil {
		t.Fatal(err)
	}
	s, err := CreateSubnetsList(set)
	if err != nil {
		t.Fatal(err)
	}
	s.Add(makesubnet("192.168.1.0/24"))
	s.Add(makesubnet("10.0.0.0/8"))

	for i := 0; i <= 255; i++ {
		ipaddr := "192.168.1." + strconv.Itoa(i)
		var h IPAddr
		if err := h.Parse(ipaddr); err != nil {
			t.Fatal(err)
		}
		if s.Check(h) != true {
			t.Fatal("SubnetsList works incorrectly at " + ipaddr)
		}
	}

	ipaddr := "192.168.2.1"
	var h IPAddr
	if err := h.Parse(ipaddr); err != nil {
		t.Fatal(err)
	}
	if s.Check(h) {
		t.Fatal("Checking not exists host")
	}

	for i := 0; i <= 255; i++ {
		ipaddr := "10.100." + strconv.Itoa(i) + "." + strconv.Itoa(255-i)
		var h IPAddr
		if err := h.Parse(ipaddr); err != nil {
			t.Fatal(err)
		}
		if s.Check(h) != true {
			t.Fatal("SubnetsList works incorrectly at " + ipaddr)
		}
	}
}

func TestSubnetsListDelete(t *testing.T) {
	set, err := createSet()
	if err != nil {
		t.Fatal(err)
	}
	s, err := CreateSubnetsList(set)
	if err != nil {
		t.Fatal(err)
	}
	s.Add(makesubnet("192.168.1.0/24"))
	deleted1, err := s.Delete(makesubnet("192.168.0.0/24"))
	if err != nil {
		t.Fatal(err)
	}
	if deleted1 {
		t.Fatal("Deleted not create subnet")
	}
	deleted2, err := s.Delete(makesubnet("192.168.1.0/24"))
	if err != nil {
		t.Fatal(err)
	}
	if !!deleted2 {
		t.Fatal("Not found added subnet")
	}
}
