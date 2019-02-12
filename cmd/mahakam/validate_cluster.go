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

type ValidateClusterOptions struct {
	Name  string
	Owner string

	ClusterAPI v1.ClusterAPI
}

var vco = &ValidateClusterOptions{}

var validateClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Validate kubernetes cluster",
	Long:  `Validate a kubernetes cluster with one command`,
	Run: func(cmd *cobra.Command, args []string) {
		if vco.Name == "" {
			exitWithHelp(cmd, "Please provide name for your cluster.")
		}

		if vco.Owner == "" {
			// Hack since we don't have login mechanism yet
			vco.Owner = config.ResourceOwnerGojek
		}

		vco.ClusterAPI = mahakamclient.GetMahakamClusterClient(os.Getenv("MAHAKAM_API_SERVER_HOST"))

		res, err := RunValidateCluster(vco)
		if err != nil {
			glog.Exit(err)
		}

		fmt.Println("Validating kubernetes cluster...")

		fmt.Println("Validating cluster nodes")
		if len(res.NodeFailures) == 0 {
			fmt.Println("PASS")
		} else {
			fmt.Println("FAIL")
			for _, nf := range res.NodeFailures {
				fmt.Printf("\t%s\n", nf)
			}
		}

		fmt.Println("Validating kubernetes system components")
		if len(res.ComponentFailures) == 0 {
			fmt.Println("PASS")
		} else {
			fmt.Println("FAIL")
			for _, cf := range res.ComponentFailures {
				fmt.Printf("\t%s\n", cf)
			}
		}

		fmt.Println("Validating pods in namespace kube-system")
		if len(res.PodFailures) == 0 {
			fmt.Println("PASS")
		} else {
			fmt.Println("FAIL")
			for _, pf := range res.PodFailures {
				fmt.Printf("\t%s\n", pf)
			}
		}
	},
}

func RunValidateCluster(vco *ValidateClusterOptions) (*models.Cluster, error) {
	req := &models.Cluster{
		Name:  swag.String(vco.Name),
		Owner: vco.Owner,
	}

	res, err := vco.ClusterAPI.ValidateCluster(clusters.NewValidateClusterParams().WithBody(req))
	if err != nil {
		return nil, fmt.Errorf("error validating cluster %s: %s", vco.Name, err)
	}

	return res.Payload, nil
}

func init() {
	// Required flags
	validateClusterCmd.Flags().StringVarP(&vco.Name, "cluster-name", "c", "", "Name for your kubernetes cluster")

	// Optional flags
	validateClusterCmd.Flags().StringVarP(&vco.Owner, "owner", "o", "", "Owner of your kubernetes cluster")

	validateCmd.AddCommand(validateClusterCmd)
}
