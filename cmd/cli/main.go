package main

import (
	"github.com/spf13/cobra"
)

func main() {

	var cmdUse = &cobra.Command{
		Use:                   "use <host>",
		Short:                 "Use host in next commands",
		Long:                  "Use (select) host as current bruteforce service.\nAll next commands will work with whese host.",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Run:                   useCommand,
	}

	var cmdReset = &cobra.Command{
		Use:                   "reset",
		Short:                 "Reset current connection to bruteforce service",
		Long:                  "Reset (disconnect) current bruteforce service.\nAfter reset command, need select new service.",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(0),
		Run:                   resetCommand,
	}

	var cmdClear = &cobra.Command{
		Use:                   "clear <login> <ip>",
		Short:                 "Clear login+host from bruteforce lists",
		Long:                  "Clear (remove) host from any blockers in service.",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(2),
		Run:                   clearCommand,
	}

	var cmdPass = &cobra.Command{
		Use:                   "pass <host|subnet>",
		Short:                 "Pass host or subnet",
		Long:                  "Add host or subnet into whitelist.\nHosts in whitelist always be passed.",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Run:                   passCommand,
	}

	var cmdUnpass = &cobra.Command{
		Use:                   "unpass <host|subnet>",
		Short:                 "Unpass host or subnet",
		Long:                  "Remove host or subnet from whitelist.\nHost will be processed within another rules.",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Run:                   unpassCommand,
	}

	var cmdBlock = &cobra.Command{
		Use:                   "block <host|subnet>",
		Short:                 "Block host or subnet",
		Long:                  "Add host or subnet into blacklist.\nHosts into blacklist always be blocked.",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Run:                   blockCommand,
	}

	var cmdUnblock = &cobra.Command{
		Use:                   "unblock <host|subnet>",
		Short:                 "Unblock host or subnet",
		Long:                  "Remove host or subnet from blacklist.\nHost will be processed within another rules.",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		Run:                   unblockCommand,
	}

	var rootCmd = &cobra.Command{Use: "cli"}
	rootCmd.AddCommand(cmdUse, cmdReset, cmdPass, cmdUnpass, cmdBlock, cmdUnblock, cmdClear)
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	err := rootCmd.Execute()
	exitOnError(err)
}
