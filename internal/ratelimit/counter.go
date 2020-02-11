package ratelimit

import (
	"time"

	"github.com/jdeal-mediamath/clockwork"
)

// Config - maximums for rates, after their, check limited
type Config struct {
	Login            int
	LoginDuration    time.Duration
	Password         int
	PasswordDuration time.Duration
	Host             int
	HostDuration     time.Duration
}

// Counter - object, which counts login attempts
type Counter struct {
	login    *Limitation
	password *Limitation
	host     *Limitation
	clock    clockwork.Clock
}

// CreateCounter - create main object
func CreateCounter(rates Config, clock clockwork.Clock) *Counter {
	var c Counter
	c.login = CreateLimitation(rates.Login, rates.LoginDuration)
	c.password = CreateLimitation(rates.Password, rates.PasswordDuration)
	c.host = CreateLimitation(rates.Host, rates.HostDuration)
	c.clock = clock
	return &c
}

// CheckAndCount - main function to count attempts and collect it in buckets
func (c *Counter) CheckAndCount(login, password, hostip string) (bool, string, error) {
	bylogin, err := c.login.Check(login)
	if err != nil {
		return false, "", err
	}
	if !bylogin {
		return false, "login rates limit", nil
	}
	bypassword, err := c.password.Check(password)
	if err != nil {
		return false, "", err
	}
	if !bypassword {
		return false, "password rates limit", nil
	}
	byhost, err := c.host.Check(hostip)
	if err != nil {
		return false, "", err
	}
	if !byhost {
		return false, "host rates limit", nil
	}
	return true, "", nil
}

// Reset - reset login+host from counter buckets
func (c *Counter) Reset(login, hostip string) (bool, error) {
	resetLogin := c.login.Reset(login)
	resetHost := c.host.Reset(hostip)
	return resetLogin || resetHost, nil
}
