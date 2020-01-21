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

func (m *subnet) Check(host ip) bool {
	return false
}
