package netsupport

import (
	"strconv"
	"testing"
)

func makesubnet(s string) Subnet {
	var snet Subnet
	snet.Parse(s)
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
