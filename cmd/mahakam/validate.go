package main

import "github.com/spf13/cobra"

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate mahakam resource",
	Long:  `Validate mahakam resource with one command`,
}

func init() {
	RootCmd.AddCommand(validateCmd)
}
