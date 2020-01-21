package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func useCommand(cmd *cobra.Command, args []string) {
	fmt.Println("Use: " + strings.Join(args, " "))
}

func resetCommand(cmd *cobra.Command, args []string) {
	fmt.Println("Reset: " + strings.Join(args, " "))
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
