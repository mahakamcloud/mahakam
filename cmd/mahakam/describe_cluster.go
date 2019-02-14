package main

import (
	"fmt"
	"os"

	"github.com/go-openapi/swag"
	"github.com/golang/glog"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/client/clusters"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	mahakamclient "github.com/mahakamcloud/mahakam/pkg/client"
	"github.com/spf13/cobra"
)

// DescribeClusterOptions represents cluster parameters
type DescribeClusterOptions struct {
	Name string
}

var dco = &DescribeClusterOptions{}

var describeClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Describe kubernetes cluster",
	Long:  `Describe a kubernetes cluster with one command`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			exitWithHelp(cmd, "Please provide name for your cluster.")
		}
		dco.Name = args[0]

		res, err := RunDescribeCluster(dco)
		if err != nil {
			glog.Exit(err)
		}

		fmt.Printf("\nName:\t%s", swag.StringValue(res.Name))
		fmt.Printf("\nOwner:\t%s", res.Owner)
		fmt.Printf("\nWorker Nodes:\t%v", res.NumNodes)
		fmt.Printf("\nStatus:\t%v", res.Status)
	},
}

// RunDescribeCluster gets description of a cluster from mahakam server
func RunDescribeCluster(dco *DescribeClusterOptions) (*models.Cluster, error) {
	c := mahakamclient.GetMahakamClient(os.Getenv("MAHAKAM_API_SERVER_HOST"))

	req := clusters.NewDescribeClustersParams()
	req.Name = swag.String(dco.Name)

	res, err := c.Clusters.DescribeClusters(req)
	if err != nil {
		return nil, fmt.Errorf("error describing cluster '%v': %v", dco, err)
	}
	return res.Payload, nil
}

func init() {
	describeCmd.AddCommand(describeClusterCmd)
}
