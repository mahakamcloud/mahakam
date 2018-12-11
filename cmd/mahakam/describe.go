package main

import "github.com/spf13/cobra"

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe mahakam resource",
	Long:  `Describe mahakam resource with one command`,
}

func init() {
	RootCmd.AddCommand(describeCmd)
}
