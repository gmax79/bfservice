package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gmax79/bfservice/internal/grpccon"
)

const timeout = time.Second * 5

var tests = []func(*grpccon.Client) error{
	testHealthCheck,
	testLimitationLogin,
	testLimitationPassword,
	testLimitationHost,
	testWhiteList,
	testBlackList,
	testLimitationLogin,
	testLimitationPassword,
	testLimitationHost,
}

type checkResult struct {
	err       error
	calls     int
	logins    map[string]int
	passwords map[string]int
	hosts     map[string]int
}

func check(conn *grpccon.Client, logins, passwords, ipaddr func() string) *checkResult {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var result checkResult
	result.logins = make(map[string]int)
	result.passwords = make(map[string]int)
	result.hosts = make(map[string]int)
	for {
		login := logins()
		password := passwords()
		ip := ipaddr()
		if login == "" || password == "" || ip == "" {
			break
		}
		result.calls++
		resp, err := conn.CheckLogin(ctx, login, password, ip)
		if resp == nil || err != nil {
			result.err = err
			return &result
		}
		if resp.Status {
			result.logins[login]++
			result.passwords[password]++
			result.hosts[ip]++
		}
	}
	return &result
}

func reset(conn *grpccon.Client, login, host string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := conn.ResetLogin(ctx, login, host)
	return err
}

func getRates(conn *grpccon.Client) (*grpccon.Rates, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	rates, err := conn.GetRates(ctx)
	if err != nil {
		return nil, err
	}
	return rates, nil
}

func calcWithLeaked(worktime time.Duration, count int, rateInteraval time.Duration) int {
	ratems := float64(count) / float64(rateInteraval.Milliseconds())
	leaked := float64(worktime.Milliseconds()) * ratems
	return count + int(leaked)
}

func testHealthCheck(conn *grpccon.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return conn.HealthCheck(ctx)
}

func testLimitationLogin(conn *grpccon.Client) error {
	fmt.Println("testLimitationLogin")
	err := reset(conn, "login", "192.168.1.1") // reset blocks for test's login and ip (to repeating tests)
	if err != nil {
		return err
	}
	rates, err := getRates(conn)
	if err != nil {
		return err
	}

	attemps := rates.LoginRate
	logins := stringGenerator(10, attemps+10, "login")
	passwords := randomString
	ip := fromConstGenerator("192.168.1.1", attemps+20)

	startTime := time.Now()
	res := check(conn, logins, passwords, ip)
	workTime := time.Now().Sub(startTime)

	testLoginsRate := res.logins["login"]
	calcLoginRate := calcWithLeaked(workTime, rates.LoginRate, rates.LoginInterval)
	fmt.Printf("limits result: calls %d, logins passed/calculated: %d/%d\n",
		res.calls, testLoginsRate, calcLoginRate)
	if calcLoginRate != testLoginsRate {
		return errors.New("testLimitationLogin failed")
	}
	fmt.Println("pass: limits as service settings")
	return res.err
}

func testLimitationPassword(conn *grpccon.Client) error {
	fmt.Println("testLimitationPassword")
	err := reset(conn, "login", "192.168.1.1") // reset blocks for test's login and ip (to repeating tests)
	if err != nil {
		return err
	}
	rates, err := getRates(conn)
	if err != nil {
		return err
	}

	randomPassword := randomString() // use random password, exclude conflicts after restart test (blocking by password)
	attemps := rates.PasswordRate
	logins := randomString
	passwords := fromConstGenerator(randomPassword, attemps+20)
	ip := fromConstGenerator("192.168.1.1", attemps+20)

	startTime := time.Now()
	res := check(conn, logins, passwords, ip)
	workTime := time.Now().Sub(startTime)

	testPasswordRate := res.passwords[randomPassword]
	calcPasswordRate := calcWithLeaked(workTime, rates.PasswordRate, rates.PasswordInterval)

	fmt.Printf("limits result: calls %d, passwords passed/calculated: %d/%d\n",
		res.calls, testPasswordRate, calcPasswordRate)
	if calcPasswordRate != testPasswordRate {
		return errors.New("testLimitationPassword failed")
	}
	fmt.Println("pass: limits as service settings")
	return res.err
}

func testLimitationHost(conn *grpccon.Client) error {
	fmt.Println("testLimitationHost")
	host := "192.168.2.1"
	err := reset(conn, "", host)
	if err != nil {
		return err
	}
	rates, err := getRates(conn)
	if err != nil {
		return err
	}

	attempts := rates.HostRate
	passwords := randomString
	logins := randomString
	ip := ipGenerator(300, "192.168.2.0", attempts+100, host)

	startTime := time.Now()
	res := check(conn, logins, passwords, ip)
	workTime := time.Now().Sub(startTime)

	calcHostRate := calcWithLeaked(workTime, rates.HostRate, rates.HostInterval)
	testHostRate := res.hosts[host]
	fmt.Printf("limits result: calls %d, ip passed/calculated %d/%d\n", res.calls, testHostRate, calcHostRate)
	if calcHostRate != testHostRate {
		return errors.New("testLimitationHost failed")
	}
	fmt.Println("pass: limits as service settings")
	return res.err
}

func testWhiteList(conn *grpccon.Client) error {
	fmt.Println("testWhiteList")
	rates, err := getRates(conn)
	if err != nil {
		return err
	}

	const host = "192.168.3.1"
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	fmt.Println("Add into whitelist", host)
	result, err := conn.AddWhiteList(ctx, host)
	printResult(result, err)
	if err != nil {
		return err
	}

	hostCalls := rates.HostRate + 100

	logins := randomString
	passwords := randomString
	ip := ipGenerator(200, "192.168.3.0", hostCalls, host)
	res := check(conn, logins, passwords, ip)

	testIPRate := res.hosts[host]
	fmt.Printf("limits result: calls %d, passed ip: %d\n", res.calls, testIPRate)
	if hostCalls != testIPRate {
		return errors.New("testWhiteList failed")
	}
	fmt.Println("pass: limits as service settings")
	return res.err
}

func testBlackList(conn *grpccon.Client) error {
	fmt.Println("testBlackList")
	rates, err := getRates(conn)
	if err != nil {
		return err
	}

	const host = "192.168.4.1"
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	fmt.Println("Add into blacklist", host)
	result, err := conn.AddBlackList(ctx, host)
	printResult(result, err)
	if err != nil {
		return err
	}

	attempts := rates.HostRate + 100
	logins := randomString
	passwords := randomString
	ip := ipGenerator(200, "192.168.4.0", attempts, host)
	res := check(conn, logins, passwords, ip)

	testIPRate := res.hosts[host]
	fmt.Printf("limits result: calls %d, passed ip: %d\n", res.calls, testIPRate)
	if testIPRate != 0 {
		return errors.New("testBlackList failed")
	}
	fmt.Println("pass: limits as service settings")
	return res.err
}
