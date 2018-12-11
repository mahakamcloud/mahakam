package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

type CreateClusterOptions struct {
	NumNodes int
}

var cco = &CreateClusterOptions{}

var createClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Create kubernetes cluster",
	Long:  `Create a kubernetes cluster with one command`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunCreateCluster(cco); err != nil {
			fmt.Printf("RunCreateCluster error: %s", err.Error())
		}
	},
}

func RunCreateCluster(cco *CreateClusterOptions) error {
	fmt.Println("RunCreateCluster not yet implemented")
	return nil
}

func init() {
	createCmd.AddCommand(createClusterCmd)
}
