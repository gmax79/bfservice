package hostslist

import (
	"strconv"
	"testing"
)

func makesubnet(s string) subnet {
	var snet subnet
	snet.Parse(s)
	return snet
}

func TestFilterInList(t *testing.T) {
	f := CreateFilter()
	f.Add(makesubnet("192.168.1.0/24"))
	f.Add(makesubnet("10.0.0.0/8"))
	for i := 0; i <= 255; i++ {
		ipaddr := "192.168.1." + strconv.Itoa(i)
		var h ip
		if err := h.Parse(ipaddr); err != nil {
			t.Fatal(err)
		}
		if f.Check(h) != true {
			t.Fatal("Filter works incorrectly at " + ipaddr)
		}
	}
	for i := 0; i <= 255; i++ {
		ipaddr := "10.100." + strconv.Itoa(i) + "." + strconv.Itoa(255-i)
		var h ip
		if err := h.Parse(ipaddr); err != nil {
			t.Fatal(err)
		}
		if f.Check(h) != true {
			t.Fatal("Filter works incorrectly at " + ipaddr)
		}
	}
}
