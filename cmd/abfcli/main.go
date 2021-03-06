package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// simple exit, without timer in log by log.Fatal
func exitOnError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	var rootCmd = &cobra.Command{Use: os.Args[0]}
	rootCmd.AddCommand(cmdUse, cmdReset, cmdPass, cmdUnpass, cmdBlock, cmdUnblock, cmdClear, cmdCheck)
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
	err := rootCmd.Execute()
	exitOnError(err)
}

var cmdUse = &cobra.Command{
	Use:                   "use <host>",
	Short:                 "Use host in next commands",
	Long:                  "Use (select) host as current bruteforce service.\nAll next commands will work with whese host.",
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := useCommand(args[0])
		exitOnError(err)
	},
}

var cmdReset = &cobra.Command{
	Use:                   "reset",
	Short:                 "Reset current connection to bruteforce service",
	Long:                  "Reset (disconnect) current bruteforce service.\nAfter reset command, need select new service.",
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		err := resetCommand()
		exitOnError(err)
	},
}

var cmdClear = &cobra.Command{
	Use:                   "clear <login> <ip>",
	Short:                 "Clear login+host from bruteforce lists",
	Long:                  "Clear (remove) host from any blockers in service.",
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := clearCommand(args[0], args[1])
		exitOnError(err)
	},
}

var cmdPass = &cobra.Command{
	Use:                   "pass <host|subnet>",
	Short:                 "Pass host or subnet",
	Long:                  "Add host or subnet into whitelist.\nHosts in whitelist always be passed.",
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := passCommand(args[0])
		exitOnError(err)
	},
}

var cmdUnpass = &cobra.Command{
	Use:                   "unpass <host|subnet>",
	Short:                 "Unpass host or subnet",
	Long:                  "Remove host or subnet from whitelist.\nHost will be processed within another rules.",
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := unpassCommand(args[0])
		exitOnError(err)
	},
}

var cmdBlock = &cobra.Command{
	Use:                   "block <host|subnet>",
	Short:                 "Block host or subnet",
	Long:                  "Add host or subnet into blacklist.\nHosts into blacklist always be blocked.",
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := blockCommand(args[0])
		exitOnError(err)
	},
}

var cmdUnblock = &cobra.Command{
	Use:                   "unblock <host|subnet>",
	Short:                 "Unblock host or subnet",
	Long:                  "Remove host or subnet from blacklist.\nHost will be processed within another rules.",
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := unblockCommand(args[0])
		exitOnError(err)
	},
}

var cmdCheck = &cobra.Command{
	Use:                   "check <login> <password> <host>",
	Short:                 "Check state in service for login/password/host",
	Long:                  "Try login with parameters and print response from service",
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		err := checkCommand(args[0], args[1], args[2])
		exitOnError(err)
	},
}
