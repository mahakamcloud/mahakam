package main

import "github.com/spf13/cobra"

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get mahakam resource",
	Long:  `Get mahakam resource with one command`,
}

func init() {
	RootCmd.AddCommand(getCmd)
}
