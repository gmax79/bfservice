package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gmax79/bfservice/internal/grpccon"
)

const host = "localhost:9000"

func exitOnError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	var err error
	defer exitOnError(err)

	fmt.Println("Automatic tests for antibruteforce service")
	conn, err := grpccon.Connect(host)
	defer conn.Close()
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	if err = conn.HealthCheck(ctx); err != nil {
		return
	}

	var resp *grpccon.Response
	if resp, err = conn.CheckLogin(ctx, "login", "password", "100.0.0.0"); err != nil {
		return
	}
	fmt.Println(*resp)
	fmt.Println("Automatic tests finished")
}
