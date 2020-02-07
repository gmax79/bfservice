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

type testsSetProvider struct {
}

func (p *testsSetProvider) Add(item string) (bool, error) {
	return true, nil
}
func (p *testsSetProvider) Delete(item string) (bool, error) {
	return true, nil
}
func (p *testsSetProvider) Iterator() (storage.StringIterator, error) {
	return func() (string, bool) {
		return "", false
	}, nil
}

func TestSubnetsListInList(t *testing.T) {
	var p testsSetProvider
	s, err := CreateSubnetsList(&p)
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
