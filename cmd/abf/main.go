package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gmax79/bfservice/internal/buckets"
)

// RatesConfig - struct to read config with bruteforce rates
type RatesConfig struct {
	LoginRate    int `json:"login_rate"`
	PasswordRate int `json:"password_rate"`
	IPRate       int `json:"ip_rate"`
}

// Check - validate config values
func (r *RatesConfig) Check() error {
	return nil
}

func exitOnError(err error) {
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
}

func main() {
	var err error
	defer exitOnError(err)

	fmt.Println("Antibruteforce service")
	var configJSON []byte
	if configJSON, err = ioutil.ReadFile("config.json"); err != nil {
		return
	}
	ratesconfig := RatesConfig{}
	if err = json.Unmarshal(configJSON, &ratesconfig); err != nil {
		return
	}
	if err = ratesconfig.Check(); err != nil {
		return
	}
	fmt.Printf("Rates: %d/login, %d/password, %d/host\n", ratesconfig.LoginRate, ratesconfig.PasswordRate, ratesconfig.IPRate)

	filter := buckets.CreateFilter()
	host := ":9000"

	grpc, err := openGRPCServer(filter, host, nil)
	if err != nil {
		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Antibruteforce service started on:", host)

	<-stop
	grpc.Stop(context.Background())
	log.Println("Antibruteforce service stopped")
}
