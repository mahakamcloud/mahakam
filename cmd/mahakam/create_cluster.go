package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

type CreateClusterOptions struct {
	Name     string
	NumNodes int
}

var cco = &CreateClusterOptions{}

var createClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Create kubernetes cluster",
	Long:  `Create a kubernetes cluster with one command`,
	Run: func(cmd *cobra.Command, args []string) {
		if cco.Name == "" {
			exitWithHelp(cmd, "Please provide name for your cluster.")
		}

		if err := RunCreateCluster(cco); err != nil {
			fmt.Printf("RunCreateCluster error: %s", err.Error())
		}
	},
}

func RunCreateCluster(cco *CreateClusterOptions) error {
	fmt.Printf("+%v", cco)

	return nil
}

func init() {
	// Required flags
	createClusterCmd.Flags().StringVarP(&cco.Name, "cluster-name", "c", "", "Name for your kubernetes cluster")

	// Optional flags
	createClusterCmd.Flags().IntVarP(&cco.NumNodes, "num-nodes", "n", 1, "Number of worker nodes you want kubernetes cluster to run")

	createCmd.AddCommand(createClusterCmd)
}
