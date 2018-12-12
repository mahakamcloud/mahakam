package main

import (
	"fmt"

	"github.com/go-openapi/swag"

	"github.com/mahakamcloud/mahakam/pkg/api/v1/client/clusters"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
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
			fmt.Printf("create cluster error: %s", err.Error())
		}
	},
}

func RunCreateCluster(cco *CreateClusterOptions) error {
	c := GetClusterClient()
	req := &models.Cluster{
		Name:     swag.String(cco.Name),
		NumNodes: int64(cco.NumNodes),
	}

	res, err := c.Clusters.CreateCluster(clusters.NewCreateClusterParams().WithBody(req))
	if err != nil {
		return fmt.Errorf("error creating cluster '%v': '%v'", cco, err)
	}

	fmt.Printf("successfully created cluster: '%v'", res)
	return nil
}

func init() {
	// Required flags
	createClusterCmd.Flags().StringVarP(&cco.Name, "cluster-name", "c", "", "Name for your kubernetes cluster")

	// Optional flags
	createClusterCmd.Flags().IntVarP(&cco.NumNodes, "num-nodes", "n", 1, "Number of worker nodes you want kubernetes cluster to run")

	createCmd.AddCommand(createClusterCmd)
}
