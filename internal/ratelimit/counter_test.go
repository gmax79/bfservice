package ratelimit

import (
	"testing"
	"time"

	"github.com/jdeal-mediamath/clockwork"
)

func TestCounter(t *testing.T) {
	cfg := Config{
		Login:            5,
		LoginDuration:    time.Second,
		Password:         10,
		PasswordDuration: time.Second,
		Host:             15,
		HostDuration:     time.Second,
	}
	clock := clockwork.NewFakeClock()
	_, err := CreateCounter(cfg, clock)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Host = 0
	_, err = CreateCounter(cfg, clock)
	if err == nil {
		t.Fatal("Error with invalid Host rate not recognized")
	}
	cfg.Host = 15
	cfg.PasswordDuration = 0
	_, err = CreateCounter(cfg, clock)
	if err == nil {
		t.Fatal("Error with invalid PasswordDuration rate not recognized")
	}
	cfg.PasswordDuration = time.Second
	cfg.Login = 0
	_, err = CreateCounter(cfg, clock)
	if err == nil {
		t.Fatal("Error with invalid Login rate not recognized")
	}
}

func TestCounterReset(t *testing.T) {
	cfg := Config{
		Login:            5,
		LoginDuration:    time.Second,
		Password:         10,
		PasswordDuration: time.Second,
		Host:             15,
		HostDuration:     time.Second,
	}
	clock := clockwork.NewFakeClock()
	c, err := CreateCounter(cfg, clock)
	if err != nil {
		t.Fatal(err)
	}
	res, _ := c.CheckAndCount("login", "password", "127.0.0.1")
	if !res {
		t.Fatal("Check must be calculated")
	}
	r1 := c.Reset("login", "130.0.0.1")
	if r1 != true {
		t.Fatal("r1==false? login is exist in limitation")
	}
	r2 := c.Reset("login2", "127.0.0.1")
	if r2 != true {
		t.Fatal("r2==false? host is exist in limitation")
	}
}

func TestCounterRatesLimits(t *testing.T) {
	cfg := Config{
		Login:            2,
		LoginDuration:    time.Second,
		Password:         3,
		PasswordDuration: time.Second,
		Host:             4,
		HostDuration:     time.Second,
	}
	clock := clockwork.NewFakeClock()
	c, err := CreateCounter(cfg, clock)
	if err != nil {
		t.Fatal(err)
	}
	login := "login"
	password := "password"
	host := "127.0.0.1"
	added1, _ := c.CheckAndCount(login, password, host)
	added2, _ := c.CheckAndCount(login, password, host)
	added3, _ := c.CheckAndCount(login+"a", password, host)
	added4, _ := c.CheckAndCount(login+"b", password+"a", host)

	if !added1 || !added2 || !added3 || !added4 {
		t.Fatal("limits not collected")
	}
	added5, _ := c.CheckAndCount(login, password+"c", host+"1")
	added6, _ := c.CheckAndCount(login+"c", password, host+"2")
	added7, _ := c.CheckAndCount(login+"d", password+"d", host)

	if added5 || added6 || added7 {
		t.Fatal("limits all collected, cant add new")
	}
}
