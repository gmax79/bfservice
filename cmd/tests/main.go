package main

import (
	"log"
	"os"

	"github.com/gmax79/bfservice/internal/grpccon"
)

const defaultHost = "localhost:9000"

func main() {
	log.Println("Autotests for antibruteforce service")
	if len(os.Args) > 2 {
		log.Fatal("Only one parameter acceptable (service host address)")
	}
	host := defaultHost
	if len(os.Args) == 2 {
		host = os.Args[1]
	}
	log.Println("Testing service:", host)
	if err := runTests(host); err != nil {
		log.Fatal(err)
	}
	log.Println("Autotests successfully finished")
}

func printResult(r *grpccon.Response, err error) {
	reason := ""
	if r != nil {
		reason = r.Reason
	}

	if err != nil {
		log.Println("grpc server response: error", err.Error(), reason)
	} else {
		if r.Status {
			log.Println("grpc server response: success", reason)
		} else {
			log.Println("grpc server response: failed", reason)
		}
	}
}

func runTests(host string) (err error) {
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
