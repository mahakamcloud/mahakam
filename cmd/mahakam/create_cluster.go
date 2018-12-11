package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

type CreateOptions struct {
	NumNodes int
}

var co = &CreateOptions{}

var createClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Create kubernetes cluster",
	Long:  `Create a kubernetes cluster with one command`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunCreate(co); err != nil {
			fmt.Printf("RunCreate error: %s", err.Error())
		}
	},
}

func RunCreate(co *CreateOptions) error {
	fmt.Println("RunCreate not yet implemented")
	return nil
}

func init() {
	createCmd.AddCommand(createClusterCmd)
}
