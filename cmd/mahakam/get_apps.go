package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

type GetAppsOptions struct{}

var gao = &GetAppsOptions{}

var getAppsCmd = &cobra.Command{
	Use:   "apps",
	Short: "Get list of applications",
	Long:  `Get list of applications deployed into kubernetes cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunGetApps(gao); err != nil {
			fmt.Printf("RunCreateCluster error: %s", err.Error())
		}
	},
}

func RunGetApps(gao *GetAppsOptions) error {
	fmt.Println("RunGetApps not yet implemented")
	return nil
}

func init() {
	createCmd.AddCommand(getAppsCmd)
}
