package ratelimit

import (
	"fmt"
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
}

// CreateCounter - create main object
func CreateCounter(rates Config, clock clockwork.Clock) (*Counter, error) {
	var err error
	loginBucketsFactory, err := CreateBucketsFactory(rates.Login, rates.LoginDuration, clock)
	if err != nil {
		return nil, fmt.Errorf("Error in login rates. %w", err)
	}
	passwordBucketsFactory, err := CreateBucketsFactory(rates.Password, rates.PasswordDuration, clock)
	if err != nil {
		return nil, fmt.Errorf("Error in password rates. %w", err)
	}
	hostBucketsFactory, err := CreateBucketsFactory(rates.Host, rates.HostDuration, clock)
	if err != nil {
		return nil, fmt.Errorf("Error in host rates. %w", err)
	}
	var c Counter
	c.login = CreateLimitation(loginBucketsFactory)
	c.password = CreateLimitation(passwordBucketsFactory)
	c.host = CreateLimitation(hostBucketsFactory)
	return &c, nil
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
