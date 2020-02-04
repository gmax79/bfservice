package main

import (
	"encoding/json"
	"errors"
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
	LoginRate     int    `json:"login_rate"`
	PasswordRate  int    `json:"password_rate"`
	IPRate        int    `json:"ip_rate"`
	Host          string `json:"host"`
	RedisHost     string `json:"redis_host"`
	RedisPassword string `json:"redis_password"`
	RedisDB       int    `json:"redis_db"`
}

// Check - validate config values
func (r *RatesAndHostConfig) Check() error {
	if r.LoginRate <= 0 || r.PasswordRate <= 0 || r.IPRate <= 0 {
		return errors.New("One of rate parameters in config is absent, negative or zero")
	}
	if r.Host == "" {
		return errors.New("Host parameter not declared")
	}
	if r.RedisHost == "" {
		return errors.New("Redis host parameter not declared")
	}
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
	logger.Info("Redis", zap.String("host", config.RedisHost), zap.Int("db", config.RedisDB))
	grpc, err := openGRPCServer(config, logger)
	if err != nil {
		exitOnError(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	logger.Info("Antibruteforce service started", zap.String("host", config.Host))
	<-stop
	grpc.Stop()
	logger.Info("Antibruteforce service stopped")
}
