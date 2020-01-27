package main

import (
	"context"
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
		testWhiteLists,
	}
	for _, t := range tests {
		if err = t(conn); err != nil {
			return
		}
	}
	return nil
}

func testHealthCheck(conn *grpccon.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	return conn.HealthCheck(ctx)
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
