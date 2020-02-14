package netsupport

import "testing"

func TestSimpleIP(t *testing.T) {
	var s Subnet
	if err := s.Parse("1.2.3.4"); err != nil {
		t.Fatal(err)
	}
	if s.mask != 32 {
		t.Fatal("Invalid initialization subnet by ip without mask")
	}
}

func TestParseSubnet(t *testing.T) {
	var s Subnet
	err := s.Parse("100.110.120.130/24")
	if err != nil {
		t.Fatal(err)
	}
	if s.address != Packip([4]byte{100, 110, 120, 130}) {
		t.Fatal("Parsed ip not equal template")
	}
	if s.mask != 24 {
		t.Fatal("Parsed mask not equal template")
	}
	err = s.Parse("1.2.3.4/32")
	if err != nil {
		t.Fatal(err)
	}
	err = s.Parse("10.100.1000.0/24")
	if err == nil {
		t.Fatal("Must be error #1")
	}
}

func TestSubnetString(t *testing.T) {
	var s Subnet
	err := s.Parse("1.2.3.4/5")
	if err != nil {
		t.Fatal(err)
	}
	if s.String() != "1.2.3.4/5" {
		t.Fatal("subnet incorrect converted into string")
	}
}

func TestInvalidSubnetString(t *testing.T) {
	var s Subnet
	err := s.Parse("1.2.3.4/33")
	if err == nil {
		t.Fatal("Cant be mask 33")
	}
	err = s.Parse("1.2.3.4/10/12")
	if err == nil {
		t.Fatal("Subnet must be with one network mask")
	}
}
