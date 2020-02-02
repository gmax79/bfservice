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
	//testAddWhiteList,
	//testWhiteLists,
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
			result.hosts[host]++
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

func testLimitationLogin(conn *grpccon.Client) (err error) {

	reset(conn, "login", "192.168.1.1")

	logins := stringGenerator(150, 50, "login")
	passwords := fromConstGenerator("password", 200)
	ip := fromConstGenerator("192.16.1.1", 200)
	res := check(conn, logins, passwords, ip)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	rates, err := conn.GetState(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("calls %d, passed login 'login': %d, password 'password': %d\n", res.calls, res.logins["login"], res.passwords["password"])
	testLoginsRate := res.logins["login"]
	testPasswordRate := res.passwords["password"]
	if rates.LoginRate != testLoginsRate || rates.PasswordRate != testPasswordRate {
		return errors.New("test rates not equal, something wrong")
	}
	return res.err
}

func testAddWhiteList(conn *grpccon.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := conn.AddWhiteList(ctx, "192.168.0.0")
	printResult(res, err)
	return err
}

func testWhiteLists(conn *grpccon.Client) (err error) {
	var resp *grpccon.Response
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if resp, err = conn.CheckLogin(ctx, "login", "password", "100.0.0.0"); err != nil {
		return
	}
	printResult(resp, err)
	return nil
}
