package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func getBinaryName() string {
	return os.Args[0]
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
		return "", fmt.Errorf("service not selected, type '%s use <abf ip[:port]>' command first", getBinaryName())
	}
	if info.IsDir() {
		return "", fmt.Errorf("fatal error, found dir with settings file path")
	}
	data, err := ioutil.ReadFile(settingsFilePath())
	if err != nil {
		return "", err
	}
	path := string(data)
	if path == "" {
		return "", fmt.Errorf("service not set, type '%s use <afb ip[:port]>' first", getBinaryName())
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
