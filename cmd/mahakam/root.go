package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "mahakam",
	Short: "application cloud platform",
	Long:  `Simple application cloud platform on Kubernetes`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func exitWithHelp(cmd *cobra.Command, err string) {
	fmt.Fprintln(os.Stderr, err)
	cmd.Help()
	os.Exit(1)
}
