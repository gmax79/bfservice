package hostslist

import (
	"errors"
	"strconv"
	"strings"
)

type ip [4]byte

func octToByte(oct string) (byte, bool) {
	val, err := strconv.Atoi(oct)
	if err != nil || val < 0 || val > 255 {
		return 0, false
	}
	return byte(val), true
}

func (v *ip) Parse(host string) error {
	var ipvalues ip
	var ok bool
	octets := strings.Split(host, ".")
	if len(octets) == 4 {
		for i, o := range octets {
			ipvalues[i], ok = octToByte(o)
			if !ok {
				break
			}
		}
	}
	if !ok {
		return errors.New("It is not correct ipv4 address: " + host)
	}
	return nil
}
