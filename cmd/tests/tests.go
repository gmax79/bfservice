package main

import (
	"context"
	"log"
	"strings"
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

func check(conn *grpccon.Client, logins, passwords, ipaddr func() string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		login := logins()
		password := passwords()
		ip := ipaddr()
		if login == "" || password == "" || ip == "" {
			break
		}
		resp, err := conn.CheckLogin(ctx, login, password, ip)
		time.Sleep(time.Millisecond * 5)
		if resp == nil {
			log.Println("Error, null response", err)
			return err
		}
		report := []string{login, password, ip, resp.Reason}
		resp.Reason = strings.Join(report, ",")
		printResult(resp, err)
		if err != nil {
			return err
		}
	}
	return nil
}

func testHealthCheck(conn *grpccon.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return conn.HealthCheck(ctx)
}

func testLimitationLogin(conn *grpccon.Client) (err error) {
	logins := stringGenerator(150, 50, "login")
	passwords := fromConstGenerator("password", 200)
	ip := fromConstGenerator("192.16.1.1", 200)
	err = check(conn, logins, passwords, ip)
	return nil
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
