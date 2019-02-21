package main

import (
	"fmt"
	"os"

	"github.com/go-openapi/swag"
	"github.com/golang/glog"

	"github.com/mahakamcloud/mahakam/pkg/api/v1"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/client/clusters"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	mahakamclient "github.com/mahakamcloud/mahakam/pkg/client"
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/spf13/cobra"
)

// CreateClusterOptions represents create cluster parameters
type CreateClusterOptions struct {
	Name     string
	Owner    string
	NumNodes int
	NodeSize string

	ClusterAPI v1.ClusterAPI
}

var cco = &CreateClusterOptions{}

var createClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Create kubernetes cluster",
	Long:  `Create a kubernetes cluster with one command`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			exitWithHelp(cmd, "Hello world")
		}
		cco.Name = args[0]

		if cco.Owner == "" {
			// Hack since we don't have login mechanism yet
			cco.Owner = config.ResourceOwnerGojek
		}

		cco.ClusterAPI = mahakamclient.GetMahakamClusterClient(os.Getenv("MAHAKAM_API_SERVER_HOST"))

		res, err := RunCreateCluster(cco)
		if err != nil {
			glog.Exit(err)
		}

		fmt.Println("Creating kubernetes cluster...")
		fmt.Printf("\nName:\t%s", swag.StringValue(res.Name))
		fmt.Printf("\nWorker Nodes:\t%v", res.NumNodes)
		fmt.Printf("\nStatus:\t%v", res.Status)
		fmt.Printf("\n\nUse 'mahakam describe cluster %s' to monitor the state of your cluster", swag.StringValue(res.Name))
	},
}

// RunCreateCluster requests cluster creation from mahakam server
func RunCreateCluster(cco *CreateClusterOptions) (*models.Cluster, error) {
	req := &models.Cluster{
		Name:     swag.String(cco.Name),
		Owner:    cco.Owner,
		NumNodes: int64(cco.NumNodes),
		NodeSize: swag.String(cco.NodeSize),
	}

	res, err := cco.ClusterAPI.CreateCluster(clusters.NewCreateClusterParams().WithBody(req))
	if err != nil {
		return nil, fmt.Errorf("error creating cluster '%v': '%v'", cco, err)
	}

	return res.Payload, nil
}

func init() {
	// Optional flags
	createClusterCmd.Flags().StringVarP(&cco.Owner, "owner", "o", "", "Owner of your kubernetes cluster")
	createClusterCmd.Flags().IntVarP(&cco.NumNodes, "num-nodes", "n", 1, "Number of worker nodes you want kubernetes cluster to run")
	createClusterCmd.Flags().StringVarP(&cco.NodeSize, "node-size", "s", "", "Configuration of worker nodes you want kubernetes cluster to run")

	createCmd.AddCommand(createClusterCmd)
}
