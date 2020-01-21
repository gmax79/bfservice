package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gmax79/antibf/internal/grpccon"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

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
	conn, err := grpccon.Connect(host)
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
	fmt.Println("Clear: " + strings.Join(args, " "))
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
