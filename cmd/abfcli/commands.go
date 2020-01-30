package main

import (
	"fmt"
	"time"

	"github.com/gmax79/bfservice/internal/grpccon"

	"golang.org/x/net/context"
)

const connectTimeout = time.Second * 2

type connector struct {
	client *grpccon.Client
	ctx    context.Context
	close  func()
}

func getConnector(host string) (*connector, error) {
	var err error
	if host == "" {
		host, err = getServiceHost()
		if err != nil {
			return nil, err
		}
	}
	var c connector
	c.client, err = grpccon.Connect(host)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	c.ctx = ctx
	c.close = func() {
		cancel()
		c.client.Close()
	}
	return &c, nil
}

func useCommand(host string) (err error) {
	fmt.Println("Use:", host)
	var conn *connector
	if conn, err = getConnector(host); err != nil {
		return err
	}
	defer conn.close()
	if err = conn.client.HealthCheck(conn.ctx); err != nil {
		return err
	}
	if err = saveServiceHost(host); err != nil {
		return err
	}
	fmt.Println("Successfully selected")
	return nil
}

func resetCommand() error {
	err := removeServiceHost()
	if err == nil {
		fmt.Println("Successfully reset")
	}
	return err
}

func clearCommand(login, ip string) (err error) {
	fmt.Printf("Clear: '%s' with ip: '%s'\n", login, ip)

	var conn *connector
	if conn, err = getConnector(""); err != nil {
		return
	}
	defer conn.close()

	var resp *grpccon.Response
	if resp, err = conn.client.ResetLogin(conn.ctx, login, ip); err != nil {
		return
	}
	if resp.Status {
		fmt.Printf("Successfully cleared, %s\n", resp.Reason)
	} else {
		fmt.Printf("Not cleared, %s\n", resp.Reason)
	}
	return nil
}

func passCommand(snet string) (err error) {
	fmt.Printf("Adding into whitelist (pass): '%s'\n", snet)

	var conn *connector
	if conn, err = getConnector(""); err != nil {
		return
	}

	defer conn.close()

	var resp *grpccon.Response
	if resp, err = conn.client.AddWhiteList(conn.ctx, snet); err != nil {
		return
	}
	if resp.Status {
		fmt.Printf("Successfully added in whitelist, %s\n", resp.Reason)
	} else {
		fmt.Printf("Not added in whitelist, %s\n", resp.Reason)
	}
	return nil
}

func unpassCommand(snet string) (err error) {
	fmt.Printf("Removing from whitelist (unpass): '%s'\n", snet)

	var conn *connector
	if conn, err = getConnector(""); err != nil {
		return
	}
	defer conn.close()

	var resp *grpccon.Response
	if resp, err = conn.client.DeleteWhiteList(conn.ctx, snet); err != nil {
		return
	}
	if resp.Status {
		fmt.Printf("Successfully removed from whitelist, %s\n", resp.Reason)
	} else {
		fmt.Printf("Not removed from whitelist, %s\n", resp.Reason)
	}
	return nil
}

func blockCommand(snet string) (err error) {
	fmt.Printf("Adding into blacklist (block): '%s'\n", snet)

	var conn *connector
	if conn, err = getConnector(""); err != nil {
		return
	}
	defer conn.close()

	var resp *grpccon.Response
	if resp, err = conn.client.AddBlackList(conn.ctx, snet); err != nil {
		return
	}
	if resp.Status {
		fmt.Printf("Successfully added in blacklist, %s\n", resp.Reason)
	} else {
		fmt.Printf("Not added in blacklist, %s\n", resp.Reason)
	}
	return nil
}

func unblockCommand(snet string) (err error) {
	fmt.Printf("Removing from blacklist (unblock): '%s'\n", snet)

	var conn *connector
	if conn, err = getConnector(""); err != nil {
		return
	}
	defer conn.close()

	var resp *grpccon.Response
	if resp, err = conn.client.DeleteBlackList(conn.ctx, snet); err != nil {
		return
	}
	if resp.Status {
		fmt.Printf("Successfully removed from blacklist, %s\n", resp.Reason)
	} else {
		fmt.Printf("Not removed from blacklist, %s\n", resp.Reason)
	}
	return nil
}
