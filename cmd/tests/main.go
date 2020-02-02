package main

import (
	"log"

	"github.com/gmax79/bfservice/internal/grpccon"
)

const host = "localhost:9000"

func main() {
	log.Println("Autotests for antibruteforce service")
	if err := runTests(); err != nil {
		log.Fatal(err)
	}
	log.Println("Autotests successfully finished")
}

func printResult(r *grpccon.Response, err error) {
	if err != nil {
		log.Println("Error", err.Error(), r.Reason)
	} else {
		log.Println("Ok", r.Reason)
	}
}

func runTests() (err error) {
	var conn *grpccon.Client
	conn, err = grpccon.Connect(host)
	if err != nil {
		return
	}
	defer conn.Close()
	tests := []func(*grpccon.Client) error{
		testHealthCheck,
		testAddWhiteList,
		testWhiteLists,
	}
	for _, t := range tests {
		if err = t(conn); err != nil {
			return
		}
	}
	return nil
}
