package main

import "github.com/spf13/cobra"

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create mahakam resource",
	Long:  `Create mahakam resource with one command`,
}

func init() {
	RootCmd.AddCommand(createCmd)
}
