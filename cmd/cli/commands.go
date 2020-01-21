package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/spf13/cobra"
)

func exitOnError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func useCommand(cmd *cobra.Command, args []string) {
	fmt.Println("Use: " + strings.Join(args, " "))
	host := args[0]
	conn, err := grpcConnect(host)
	exitOnError(err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err = conn.HealthCheck(ctx)
	exitOnError(err)
	err = saveServiceHost(host)
	exitOnError(err)
	fmt.Println("Successfully selected")
}

func resetCommand(cmd *cobra.Command, args []string) {
	//fmt.Println("Reset: " + strings.Join(args, " "))
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
