package hostslist

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type subnet struct {
	address ip
	mask    int
}

func (m *subnet) Parse(subnet string) error {
	var err error
	parts := strings.Split(subnet, "/")
	if len(parts) == 2 {
		err = m.address.Parse(parts[0])
		if err == nil {
			m.mask, err = strconv.Atoi(parts[1])
			if m.mask < 0 || m.mask > 32 {
				err = errors.New("Invalid mask: " + parts[1])
			}
		}
	}
	if err != nil {
		return fmt.Errorf("It is not correct subnet address: %s %s,", subnet, err.Error())
	}
	return nil
}

func (m *subnet) String() string {
	return fmt.Sprintf("%s/%d", m.address.String(), m.mask)
}

var masks = [32]ip{}

func init() {
	masks[0] = 0
	var mask ip = 0x8000
	for i := 1; i < 32; i++ {
		masks[i] = masks[i-1] | mask
		mask = mask >> 1
	}
}

func (m *subnet) Check(host ip) bool {
	return (host & masks[m.mask]) == m.address
}
