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

	"go.uber.org/zap"

	nlog "github.com/gmax79/bfservice/internal/log"
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

func main() {
	var err error
	var configJSON []byte
	if configJSON, err = ioutil.ReadFile("config.json"); err != nil {
		log.Fatal(err)
	}
	var logger *zap.Logger
	if logger, err = nlog.CreateLogger(configJSON); err != nil {
		log.Fatal(err)
	}

	exitOnError := func(err error) {
		logger.Error("Error", zap.Error(err))
		os.Exit(1)
	}

	logger.Info("Antibruteforce service")

	ratesconfig := RatesConfig{}
	if err = json.Unmarshal(configJSON, &ratesconfig); err != nil {
		exitOnError(err)
	}
	if err = ratesconfig.Check(); err != nil {
		exitOnError(err)
	}

	fmt.Printf("Rates: %d/login, %d/password, %d/host\n", ratesconfig.LoginRate, ratesconfig.PasswordRate, ratesconfig.IPRate)
	host := ":9000"
	grpc, err := openGRPCServer(host, logger)
	if err != nil {
		exitOnError(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	logger.Info("Antibruteforce service started", zap.String("host", host))
	<-stop
	grpc.Stop(context.Background())
	logger.Info("Antibruteforce service stopped")
}
