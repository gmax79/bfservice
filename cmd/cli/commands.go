package main

import (
	"fmt"
	"os"
	"strings"
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

func useCommand(cmd *cobra.Command, args []string) {
	var err error
	defer exitOnError(err)
	host := args[0]
	fmt.Println("Use:", host)

	var conn *grpccon.Client
	conn, err = grpccon.Connect(host)
	defer conn.Close()
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	if err = conn.HealthCheck(ctx); err != nil {
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

	host, err := getServiceHost()
	if err != nil {
		return
	}

	var conn *grpccon.Client
	conn, err = grpccon.Connect(host)
	defer conn.Close()
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	var resp *grpccon.Response
	resp, err = conn.ResetLogin(ctx, login, ip)
	if err != nil {
		return
	}
	if resp.Status {
		fmt.Printf("Successfully cleared, %s\n", resp.Reason)
	} else {
		fmt.Printf("Not cleared, %s\n", resp.Reason)
	}
}

func passCommand(cmd *cobra.Command, args []string) {
	fmt.Println("Pass: " + strings.Join(args, " "))
}

func unpassCommand(cmd *cobra.Command, args []string) {
	fmt.Println("Unpass: " + strings.Join(args, " "))
}

func blockCommand(cmd *cobra.Command, args []string) {
	fmt.Println("Block: " + strings.Join(args, " "))
}

func unblockCommand(cmd *cobra.Command, args []string) {
	fmt.Println("Unblock: " + strings.Join(args, " "))
}
