package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

type CreateAppOptions struct{}

var cao = &CreateAppOptions{}

var createAppCmd = &cobra.Command{
	Use:   "app",
	Short: "Create application",
	Long:  `Create application on kubernetes cluster with one command`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunCreateApp(cao); err != nil {
			fmt.Printf("RunCreateApp error: %s", err.Error())
		}
	},
}

func RunCreateApp(cao *CreateAppOptions) error {
	fmt.Println("RunCreateApp not yet implemented")
	return nil
}

func init() {
	createCmd.AddCommand(createAppCmd)
}
