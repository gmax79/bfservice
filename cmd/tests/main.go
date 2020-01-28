package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gmax79/bfservice/internal/grpccon"
)

const host = "localhost:9000"

func main() {
	log.Println("Autotests for antibruteforce service")
	if err := runTests(); err != nil {
		log.Fatal(err)
	}
	log.Println("Autotests finished")
}

func runTests() (err error) {
	conn, err := grpccon.Connect(host)
	defer conn.Close()
	if err != nil {
		return
	}
	tests := []func(*grpccon.Client) error{
		testHealthCheck,
		testAddWhiteList,
		//testWhiteLists,
	}
	for _, t := range tests {
		if err = t(conn); err != nil {
			return
		}
	}
	return nil
}

func printResult(r *grpccon.Response, err error) {
	if err != nil {
		fmt.Println("Error", err.Error(), r.Reason)
	} else {
		fmt.Println("Ok", r.Reason)
	}
}

func testHealthCheck(conn *grpccon.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	return conn.HealthCheck(ctx)
}

func testAddWhiteList(conn *grpccon.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	res, err := conn.AddWhiteList(ctx, "192.168.0.0")
	printResult(res, err)
	return err
}

func testWhiteLists(conn *grpccon.Client) (err error) {
	var resp *grpccon.Response
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	if resp, err = conn.CheckLogin(ctx, "login", "password", "100.0.0.0"); err != nil {
		return
	}
	log.Println(*resp)
	return nil
}
