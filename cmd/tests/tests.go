package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gmax79/bfservice/internal/grpccon"
)

const timeout = time.Second * 15

var tests = []func(*grpccon.Client) error{
	testHealthCheck,
	testLimitationLoginPassword,
	testLimitationHost,
	testWhiteList,
	//testLimitationHost,
	//testLimitationLoginPassword,
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
		time.Sleep(time.Millisecond * 5)
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

func testHealthCheck(conn *grpccon.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return conn.HealthCheck(ctx)
}

func testLimitationLoginPassword(conn *grpccon.Client) error {
	fmt.Println("testLimitationLoginPassword")
	err := reset(conn, "login", "192.168.1.1") // reset blocks for test's login and ip (to repeating tests)
	if err != nil {
		return err
	}
	randomPassword := randomString() // use random password, exclude conflicts after restart test (blocking by password)
	logins := stringGenerator(150, 50, "login")
	passwords := fromConstGenerator(randomPassword, 200)
	ip := fromConstGenerator("192.168.1.1", 200)
	res := check(conn, logins, passwords, ip)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	rates, err := conn.GetState(ctx)
	if err != nil {
		return err
	}
	testLoginsRate := res.logins["login"]
	testPasswordRate := res.passwords[randomPassword]
	fmt.Printf("limits result: calls %d, passed login 'login': %d, password '%s': %d\n",
		res.calls, testLoginsRate, randomPassword, testPasswordRate)
	if rates.LoginRate != testLoginsRate || rates.PasswordRate != testPasswordRate {
		return errors.New("testLimitationLoginPassword failed")
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
	passwords := randomString
	logins := randomString
	ip := ipGenerator(300, "192.168.2.0", 1100, host)
	res := check(conn, logins, passwords, ip)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	rates, err := conn.GetState(ctx)
	if err != nil {
		return err
	}
	testHostRate := res.hosts[host]
	fmt.Printf("limits result: calls %d, passed ip '%s': %d\n", res.calls, host, testHostRate)
	if rates.HostRate != testHostRate {
		return errors.New("testLimitationHost failed")
	}
	fmt.Println("pass: limits as service settings")

	return res.err
}

func testWhiteList(conn *grpccon.Client) error {
	fmt.Println("testWhiteList")
	const host = "192.168.3.1"
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	fmt.Println("Add into whitelist", host)
	result, err := conn.AddWhiteList(ctx, host)
	printResult(result, err)
	if err != nil {
		return err
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), timeout)
	defer cancel2()
	rates, err := conn.GetState(ctx2)
	if err != nil {
		return err
	}

	hostCalls := rates.HostRate + 100

	logins := randomString
	passwords := randomString
	ip := ipGenerator(200, "192.168.3.0", hostCalls, host)
	res := check(conn, logins, passwords, ip)

	testIPRate := res.hosts[host]
	fmt.Printf("limits result: calls %d, passed ip '%s': %d\n", res.calls, host, testIPRate)
	if hostCalls != testIPRate {
		return errors.New("testWhiteList failed")
	}
	fmt.Println("pass: limits as service settings")
	return res.err
}
