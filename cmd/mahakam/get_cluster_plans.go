package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

type GetClusterPlansOptions struct{}

var gcpo = &GetClusterPlansOptions{}

var getClusterPlansCmd = &cobra.Command{
	Use:   "cluster-plans",
	Short: "Get list of available cluster plans",
	Long:  `Get list of available cluster plans to create Kubernetes cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunGetClusterPlans(gcpo); err != nil {
			fmt.Printf("RunGetClusterPlans error: %s", err.Error())
		}
	},
}

func RunGetClusterPlans(gcpo *GetClusterPlansOptions) error {
	fmt.Println("RunGetClusterPlans not yet implemented")
	return nil
}

func init() {
	createCmd.AddCommand(getClusterPlansCmd)
}
