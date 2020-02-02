package main

import (
	"log"
	"math/rand"

	"github.com/gmax79/bfservice/internal/netsupport"
)

func fromArrayGenerator(elements []string) func() string {
	i := 0
	return func() string {
		if i < len(elements) {
			str := elements[i]
			i++
			return str
		}
		return ""
	}
}

func fromConstGenerator(element string) func() string {
	return func() string {
		return element
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString() string {
	n := rand.Intn(4) + 4
	b := make([]rune, n)
	ll := len(letters)
	for i := range b {
		b[i] = letters[rand.Intn(ll)]
	}
	return string(b)
}

// stringGenerator - genrate random+count strings, where count of good string
func stringGenerator(randomcount, count int, good string) func() string {
	return func() string {
		if randomcount == 0 && count == 0 {
			return ""
		}
		n := rand.Intn(100)
		if n%2 == 0 {
			if count > 0 {
				count--
				return good
			}
			randomcount--
			return randomString()
		}
		if randomcount > 0 {
			randomcount--
			return randomString()
		}
		count--
		return good
	}
}

// zero octects replaced by 100-250
func randomIP(mask string) string {
	var ip netsupport.IPAddr
	err := ip.Parse(mask)
	if err != nil {
		log.Fatal(err)
	}
	b := netsupport.Unpackip(ip)
	for i := 0; i < 4; i++ {
		if b[i] == 0 {
			b[i] = byte(rand.Intn(150) + 100)
		}
	}
	randomip := netsupport.Packip(b)
	return randomip.String()
}

func ipGenerator(randomcount int, mask string, count int, ip string) func() string {
	return func() string {
		if randomcount == 0 && count == 0 {
			return ""
		}
		n := rand.Intn(100)
		if n%2 == 0 {
			if count > 0 {
				count--
				return ip
			}
			randomcount--
			return randomIP(mask)
		}
		if randomcount > 0 {
			randomcount--
			return randomIP(mask)
		}
		count--
		return ip
	}
}
