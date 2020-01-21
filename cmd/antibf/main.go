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
)

// RatesConfig - struct to read config with bruteforce rates
type RatesConfig struct {
	LoginRate    int `json:"login_rate"`
	PasswordRate int `json:"password_rate"`
	IPRate       int `json:"id_rate"`
}

// Check - validate config values
func (r *RatesConfig) Check() error {
	return nil
}

func main() {
	var err error
	defer func() {
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}()

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
	fmt.Println("Rates:", ratesconfig)
	grpc, err := openGRPCConnect("localhost:9000", nil)
	if err != nil {
		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Antibruteforce service started")

	<-stop
	grpc.Stop(context.Background())
	log.Println("Antibruteforce service stopped")
}
