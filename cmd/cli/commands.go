package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gmax79/bfservice/internal/grpccon"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// simple exit, without timer by log.Fatal
func exitOnError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func grpcTimeout() time.Duration {
	return time.Second * 2
}

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
	ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout())
	c.ctx = ctx
	c.close = func() {
		cancel()
		c.client.Close()
	}
	return &c, nil
}

func useCommand(cmd *cobra.Command, args []string) {
	var err error
	defer exitOnError(err)
	host := args[0]
	fmt.Println("Use:", host)

	var conn *connector
	if conn, err = getConnector(host); err != nil {
		return
	}
	defer conn.close()
	if err = conn.client.HealthCheck(conn.ctx); err != nil {
		return
	}
	if err = saveServiceHost(host); err != nil {
		return
	}
	fmt.Println("Successfully selected")
}

func resetCommand(cmd *cobra.Command, args []string) {
	err := removeServiceHost()
	exitOnError(err)
	fmt.Println("Successfully reset")
}

func clearCommand(cmd *cobra.Command, args []string) {
	var err error
	defer exitOnError(err)

	login := args[0]
	ip := args[1]
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
}

func passCommand(cmd *cobra.Command, args []string) {
	var err error
	defer exitOnError(err)
	snet := args[0]
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
}

func unpassCommand(cmd *cobra.Command, args []string) {
	var err error
	defer exitOnError(err)
	snet := args[0]
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
}

func blockCommand(cmd *cobra.Command, args []string) {
	var err error
	defer exitOnError(err)
	snet := args[0]
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
}

func unblockCommand(cmd *cobra.Command, args []string) {
	var err error
	defer exitOnError(err)
	snet := args[0]
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
}
