package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

type GetChartsOptions struct{}

var gco = &GetChartsOptions{}

var getChartsCmd = &cobra.Command{
	Use:   "charts",
	Short: "Get list of available helm charts",
	Long:  `Get list of available helm charts stored in central chart repository`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunGetCharts(gco); err != nil {
			fmt.Printf("RunGetCharts error: %s", err.Error())
		}
	},
}

func RunGetCharts(gco *GetChartsOptions) error {
	fmt.Println("RunGetCharts not yet implemented")
	return nil
}

func init() {
	getCmd.AddCommand(getChartsCmd)
}
