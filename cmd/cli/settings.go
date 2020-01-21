package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

var currentServiceHost string

func settingsFilePath() string {
	userHomedir := os.Getenv("HOME")
	return path.Join(userHomedir, ".antibf/host")
}

func saveServiceHost(host string) error {
	path := settingsFilePath()
	dir, _ := filepath.Split(path)
	var err error
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0700)
	}
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, []byte(host), 0600)
	if err == nil {
		currentServiceHost = host
	}
	return err
}

func getServiceHost() (string, error) {
	if currentServiceHost != "" {
		return currentServiceHost, nil
	}
	data, err := ioutil.ReadFile(settingsFilePath())
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func removeServiceHost() error {
	err := os.Remove(settingsFilePath())
	if err == nil || os.IsNotExist(err) {
		currentServiceHost = ""
		return nil
	}
	return err
}
