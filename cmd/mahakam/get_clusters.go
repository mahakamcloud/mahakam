package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

type GetClustersOptions struct{}

var gclo = &GetClustersOptions{}

var getClustersCmd = &cobra.Command{
	Use:   "clusters",
	Short: "Get list of kubernetes clusters",
	Long:  `Get list of kubernetes clusters owned by you`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunGetClusters(gclo); err != nil {
			fmt.Printf("RunGetCharts error: %s", err.Error())
		}
	},
}

func RunGetClusters(gclo *GetClustersOptions) error {
	fmt.Println("RunGetClusters not yet implemented")
	return nil
}

func init() {
	getCmd.AddCommand(getClustersCmd)
}
