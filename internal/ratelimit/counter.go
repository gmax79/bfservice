package ratelimit

import (
	"time"
)

// RatesLimits - maximums for rates, after their, check limited
type RatesLimits struct {
	Login    int
	Password int
	Host     int
}

// AttemptsCounter - struct, which counts login attempts
type AttemptsCounter struct {
	login    *Limitation
	password *Limitation
	host     *Limitation
}

// CreateCounter - create main object
func CreateCounter(rates RatesLimits) *AttemptsCounter {
	var c AttemptsCounter
	c.login = CreateLimitation(rates.Login, time.Minute)
	c.password = CreateLimitation(rates.Password, time.Minute)
	c.host = CreateLimitation(rates.Host, time.Minute)
	return &c
}

// CheckAndCount - main function to count attempts and collect it in buckets
func (c *AttemptsCounter) CheckAndCount(login, password, hostip string) (bool, string, error) {
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
func (c *AttemptsCounter) Reset(login, hostip string) (bool, error) {
	resetLogin := c.login.Reset(login)
	resetHost := c.host.Reset(hostip)
	return resetLogin || resetHost, nil
}
