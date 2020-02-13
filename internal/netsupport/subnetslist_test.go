package netsupport

import (
	"log"
	"strconv"
	"testing"
)

func makesubnet(s string) Subnet {
	var snet Subnet
	err := snet.Parse(s)
	if err != nil {
		log.Fatal(err)
	}
	return snet
}

func TestSubnetsListInList(t *testing.T) {
	s := CreateSubnetsList()
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
	s := CreateSubnetsList()
	s.Add(makesubnet("192.168.1.0/24"))
	deleted1 := s.Delete(makesubnet("192.168.0.0/24"))
	if deleted1 {
		t.Fatal("Deleted not create subnet")
	}
	exist := s.Exist(makesubnet("192.168.1.0/24"))
	deleted2 := s.Delete(makesubnet("192.168.1.0/24"))
	if !exist || !deleted2 {
		t.Fatal("Not found added subnet")
	}
}
