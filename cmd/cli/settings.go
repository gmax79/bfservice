package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

var binaryName string

func init() {
	binaryName = os.Args[0]
}

func settingsFilePath() string {
	userHomedir := os.Getenv("HOME")
	return path.Join(userHomedir, ".abf/host")
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
	return err
}

func getServiceHost() (string, error) {
	settingsFile := settingsFilePath()
	info, err := os.Stat(settingsFile)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("Service not selected, type '%s use <abf ip[:port]>' command first", binaryName)
	}
	if info.IsDir() {
		return "", fmt.Errorf("Fatal error, found dir with settings file path")
	}
	data, err := ioutil.ReadFile(settingsFilePath())
	if err != nil {
		return "", err
	}
	path := string(data)
	if path == "" {
		return "", fmt.Errorf("Service not set, type '%s use <afb ip[:port]>' first", binaryName)
	}
	return path, nil
}

func removeServiceHost() error {
	err := os.Remove(settingsFilePath())
	if err == nil || os.IsNotExist(err) {
		return nil
	}
	return err
}
