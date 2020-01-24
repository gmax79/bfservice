package netsupport

import (
	"errors"
	"strconv"
	"testing"
)

func TestOctoToByte(t *testing.T) {
	for i := 0; i <= 255; i++ {
		s := strconv.Itoa(i)
		b, ok := octToByte(s)
		if !ok || b != byte(i) {
			t.Fatal("OctoToByte work incorrect at " + s)
		}
	}
}

func TestPackUnpack(t *testing.T) {
	b := [4]byte{123, 140, 190, 220}
	packedip := packip(b)
	unpackedip := unpackip(packedip)
	if b != unpackedip {
		t.Fatal("Packip and unpackip working inconsistent")
	}
}

func checkip(sip string, a, b, c, d byte) error {
	var v IPAddr
	if err := v.Parse(sip); err != nil {
		return err
	}
	unpackedip := unpackip(v)
	if unpackedip[0] == a && unpackedip[1] == b && unpackedip[2] == c && unpackedip[3] == d {
		return nil
	}
	return errors.New("not correct parsed ip: " + sip)
}

func TestParseIP(t *testing.T) {
	if err := checkip("221.121.75.82", 221, 121, 75, 82); err != nil {
		t.Fatal(err)
	}
	if err := checkip("1.2.4.8", 1, 2, 4, 8); err != nil {
		t.Fatal(err)
	}
	if err := checkip("100.99.98.97", 100, 99, 98, 97); err != nil {
		t.Fatal(err)
	}
	v := packip([4]byte{10, 20, 30, 40})
	if v.String() != "10.20.30.40" {
		t.Fatal("ip incorrect converted into string")
	}
}

func TestIncorrectParseIP(t *testing.T) {
	var v IPAddr
	if err := v.Parse("105.11.1.256"); err == nil {
		t.Fatal("Must be error #1")
	}
	if err := v.Parse("105a11.1.240"); err == nil {
		t.Fatal("Must be error #2")
	}
	if err := v.Parse(""); err == nil {
		t.Fatal("Must be error #3")
	}
	if err := v.Parse("1.2.3.4.5"); err == nil {
		t.Fatal("Must be error #4")
	}
}
