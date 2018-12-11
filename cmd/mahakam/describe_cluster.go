package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

type DescribeClusterOptions struct{}

var dco = &DescribeClusterOptions{}

var describeClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Describe kubernetes cluster",
	Long:  `Describe a kubernetes cluster with one command`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunDescribeCluster(dco); err != nil {
			fmt.Printf("RunDescribeCluster error: %s", err.Error())
		}
	},
}

func RunDescribeCluster(dco *DescribeClusterOptions) error {
	fmt.Println("RunDescribeCluster not yet implemented")
	return nil
}

func init() {
	describeCmd.AddCommand(createClusterCmd)
}
