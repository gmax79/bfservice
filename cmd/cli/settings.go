package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

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
		return "", fmt.Errorf("Service not selected, type 'use <service>' first")
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
		return "", fmt.Errorf("Service not set, type 'use <service>' first")
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
