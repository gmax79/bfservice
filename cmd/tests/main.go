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
		log.Println("grpc response: error", err.Error(), r.Reason)
	} else {
		if r.Status {
			log.Println("grpc response: success", r.Reason)
		} else {
			log.Println("grpc response: failed", r.Reason)
		}
	}
}

func runTests() (err error) {
	var conn *grpccon.Client
	conn, err = grpccon.Connect(host)
	if err != nil {
		return
	}
	defer conn.Close()
	// tests - array from tests.go
	for _, t := range tests {
		if err = t(conn); err != nil {
			return
		}
	}
	return nil
}
