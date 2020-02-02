package netsupport

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// IPAddr - object to store ip address
type IPAddr uint32

func octToByte(oct string) (byte, bool) {
	val, err := strconv.Atoi(oct)
	if err != nil || val < 0 || val > 255 {
		return 0, false
	}
	return byte(val), true
}

func Packip(b [4]byte) IPAddr {
	v := binary.BigEndian.Uint32(b[:])
	return IPAddr(v)
}

func Unpackip(v IPAddr) [4]byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(v))
	var a [4]byte
	copy(a[:], b[:4])
	return a
}

func (v *IPAddr) String() string {
	b := Unpackip(*v)
	return fmt.Sprintf("%d.%d.%d.%d", b[0], b[1], b[2], b[3])
}

// Parse - convert ip in string format into object
func (v *IPAddr) Parse(host string) error {
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
		*v = Packip(ipvalues)
		return nil
	}
	return errors.New("It is not correct ipv4 address: " + host)
}
