package main

import (
	"fmt"
	"os"

	"github.com/go-openapi/swag"
	"github.com/golang/glog"
	v1 "github.com/mahakamcloud/mahakam/pkg/api/v1"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/client/clusters"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/handlers"
	"github.com/spf13/cobra"
)

type GetClustersOptions struct {
	Owner string

	ClusterAPI v1.ClusterAPI
}

var gclo = &GetClustersOptions{}

var getClustersCmd = &cobra.Command{
	Use:   "clusters",
	Short: "Get list of kubernetes clusters",
	Long:  `Get list of kubernetes clusters owned by you`,
	Run: func(cmd *cobra.Command, args []string) {
		if gclo.Owner == "" {
			// TODO: ResourceOwnerGojek should not be default owner
			gclo.Owner = config.ResourceOwnerGojek
		}

		gclo.ClusterAPI = handlers.GetMahakamClusterClient(os.Getenv("MAHAKAM_API_SERVER_HOST"))
		clusters, err := RunGetClusters(gclo)
		if err != nil {
			glog.Exit(err)
		}

		// TODO: Move cluster list formatting in separate method
		fmt.Println("Name\t\tSize\tNodes\tStatus")
		for _, v := range clusters {
			fmt.Printf("%s\t%s\t%d\t%s\n", *v.Name, v.ClusterPlan, v.NumNodes, v.Status)
		}
	},
}

func RunGetClusters(gclo *GetClustersOptions) ([]*models.Cluster, error) {
	owner := swag.String(gclo.Owner)
	res, err := gclo.ClusterAPI.GetClusters(clusters.NewGetClustersParams().WithOwner(owner))

	if err != nil {
		return nil, fmt.Errorf("Error getting cluster '%v': '%v'", gclo, err)
	}

	return res.Payload, nil
}

func init() {
	getClustersCmd.Flags().StringVarP(&gclo.Owner, "owner", "o", "", "Owner of your Kubernetes Cluster")
	getCmd.AddCommand(getClustersCmd)
}
