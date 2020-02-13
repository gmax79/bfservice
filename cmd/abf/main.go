package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	nlog "github.com/gmax79/bfservice/internal/log"
)

// RatesAndHostConfig - struct to read config with bruteforce rates and host parameters
type RatesAndHostConfig struct {
	LoginRate    int    `json:"login_rate"`
	PasswordRate int    `json:"password_rate"`
	IPRate       int    `json:"ip_rate"`
	Host         string `json:"host"`
}

// Check - validate config values
func (r *RatesAndHostConfig) Check() error {
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

	var config RatesAndHostConfig
	if err = json.Unmarshal(configJSON, &config); err != nil {
		exitOnError(err)
	}
	if err = config.Check(); err != nil {
		exitOnError(err)
	}

	logger.Info("Rates", zap.Int("login", config.LoginRate),
		zap.Int("password", config.PasswordRate), zap.Int("host", config.IPRate))
	grpc, err := openGRPCServer(config, logger)
	if err != nil {
		exitOnError(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	logger.Info("Antibruteforce service started", zap.String("host", config.Host))
	<-stop
	logger.Info("Stopping abf service")
	grpc.Stop()
	logger.Info("Antibruteforce service stopped")
}
