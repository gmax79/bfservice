package hostslist

import "testing"

func TestParseSubnet(t *testing.T) {
	var s subnet
	err := s.Parse("100.110.120.130/24")
	if err != nil {
		t.Fatal(err)
	}

}
