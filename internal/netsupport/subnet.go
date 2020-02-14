package netsupport

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Subnet - struct to store informaion about network
type Subnet struct {
	address IPAddr
	mask    int
}

// Parse - convert subnet in string format into object
func (m *Subnet) Parse(subnet string) error {
	var err error
	parts := strings.Split(subnet, "/")
	switch len(parts) {
	case 1:
		err = m.address.Parse(parts[0])
		if err == nil {
			m.mask = 32
		}
	case 2:
		err = m.address.Parse(parts[0])
		if err == nil {
			m.mask, err = strconv.Atoi(parts[1])
			if m.mask < 0 || m.mask > 32 {
				err = errors.New("Invalid mask: " + parts[1])
			}
		}
	default:
		return fmt.Errorf("incorrect subnet address: %s", subnet)
	}
	if err != nil {
		return fmt.Errorf("incorrect subnet address: %s %s,", subnet, err.Error())
	}
	return nil
}

func (m *Subnet) String() string {
	return fmt.Sprintf("%s/%d", m.address.String(), m.mask)
}

var masks = [33]IPAddr{}

func init() {
	masks[0] = 0
	var mask IPAddr = 0x80000000
	for i := 1; i <= 32; i++ {
		masks[i] = masks[i-1] | mask
		mask >>= 1
	}
}

// Check - check if host in subnet
func (m *Subnet) Check(host IPAddr) bool {
	masked := host & masks[m.mask]
	return masked == m.address
}
