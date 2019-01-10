package main

import (
	"fmt"
	"os"

	"github.com/go-openapi/swag"
	"github.com/golang/glog"

	"github.com/mahakamcloud/mahakam/pkg/api/v1"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/client/clusters"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/handlers"
	"github.com/spf13/cobra"
)

type CreateClusterOptions struct {
	Name     string
	Owner    string
	NumNodes int

	ClusterAPI v1.ClusterAPI
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

		if cao.Owner == "" {
			// Hack since we don't have login mechanism yet
			cco.Owner = config.ResourceOwnerGojek
		}

		cco.ClusterAPI = handlers.GetMahakamClusterClient(os.Getenv("MAHAKAM_API_SERVER_HOST"))

		res, err := RunCreateCluster(cco)
		if err != nil {
			glog.Exit(err)
		}

		fmt.Println("Creating kubernetes cluster...")
		fmt.Printf("\nName:\t%s", swag.StringValue(res.Name))
		fmt.Printf("\nCluster Plan:\t%s", res.ClusterPlan)
		fmt.Printf("\nWorker Nodes:\t%v", res.NumNodes)
		fmt.Printf("\nStatus:\t%v", res.Status)
		fmt.Printf("\n\nUse 'mahakam describe cluster %s' to monitor the state of your cluster", swag.StringValue(res.Name))
	},
}

func RunCreateCluster(cco *CreateClusterOptions) (*models.Cluster, error) {
	req := &models.Cluster{
		Name:     swag.String(cco.Name),
		Owner:    cco.Owner,
		NumNodes: int64(cco.NumNodes),
	}

	res, err := cco.ClusterAPI.CreateCluster(clusters.NewCreateClusterParams().WithBody(req))
	if err != nil {
		return nil, fmt.Errorf("error creating cluster '%v': '%v'", cco, err)
	}

	return res.Payload, nil
}

func init() {
	// Required flags
	createClusterCmd.Flags().StringVarP(&cco.Name, "cluster-name", "c", "", "Name for your kubernetes cluster")

	// Optional flags
	createClusterCmd.Flags().StringVarP(&cco.Owner, "owner", "o", "", "Owner of your kubernetes cluster")
	createClusterCmd.Flags().IntVarP(&cco.NumNodes, "num-nodes", "n", 1, "Number of worker nodes you want kubernetes cluster to run")

	createCmd.AddCommand(createClusterCmd)
}
