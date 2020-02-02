package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gmax79/bfservice/internal/grpccon"
)

const timeout = time.Second * 2

var tests = []func(*grpccon.Client) error{
	testHealthCheck,
	testLimitationLogin,
	//testAddWhiteList,
	//testWhiteLists,
}

func check(conn *grpccon.Client, logins, passwords, ipaddr func() string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for login := logins(); login != ""; login = logins() {
		for password := passwords(); password != ""; password = passwords() {
			for ip := ipaddr(); ip != ""; ip = ipaddr() {
				resp, err := conn.CheckLogin(ctx, login, password, ip)
				printResult(resp, err)
				if err != nil {
					return err
				}
			}
		}
	}
	for login := logins(); login != ""; login = logins() {
		fmt.Println(login)
	}

	return nil
}

func testHealthCheck(conn *grpccon.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return conn.HealthCheck(ctx)
}

func testAddWhiteList(conn *grpccon.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := conn.AddWhiteList(ctx, "192.168.0.0")
	printResult(res, err)
	return err
}

func testLimitationLogin(conn *grpccon.Client) (err error) {
	logins := stringGenerator(15, 5, "login")
	passwords := fromConstGenerator("password", 20)
	ip := fromConstGenerator("192.16.1.1", 20)
	err = check(conn, logins, passwords, ip)
	return nil
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
