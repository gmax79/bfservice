package hostslist

import "testing"

func makesubnet(s string) subnet {
	var snet subnet
	snet.Parse(s)
	return snet
}

func TestFilter(t *testing.T) {

	f := CreateFilter()

	f.Add(s)
}
