package hostslist

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type ip uint32

func octToByte(oct string) (byte, bool) {
	val, err := strconv.Atoi(oct)
	if err != nil || val < 0 || val > 255 {
		return 0, false
	}
	return byte(val), true
}

func packip(b [4]byte) ip {
	v := binary.BigEndian.Uint32(b[:])
	return ip(v)
}

func unpackip(v ip) [4]byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(v))
	var a [4]byte
	copy(a[:], b[:4])
	return a
}

func (v *ip) String() string {
	b := unpackip(*v)
	return fmt.Sprintf("%d.%d.%d.%d", b[0], b[1], b[2], b[3])
}

func (v *ip) Parse(host string) error {
	var ipvalues [4]byte
	var ok bool
	octets := strings.Split(host, ".")
	if len(octets) == 4 {
		for i, o := range octets {
			if ipvalues[i], ok = octToByte(o); !ok {
				break
			}
		}
	}
	if ok {
		*v = packip(ipvalues)
		return nil
	}
	return errors.New("It is not correct ipv4 address: " + host)
}
