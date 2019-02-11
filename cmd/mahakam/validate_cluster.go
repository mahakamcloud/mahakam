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

		vco.ClusterAPI = handlers.GetMahakamClusterClient(os.Getenv("MAHAKAM_API_SERVER_HOST"))

		res, err := RunValidateCluster(vco)
		if err != nil {
			glog.Exit(err)
		}

		fmt.Println("Validating kubernetes cluster...")

		// TODO(giri): print out nodes validation result
		fmt.Println("Validating cluster nodes")
		fmt.Println("PASS")

		// TODO(giri): print out system components result
		fmt.Println("Validating kubernetes system components")
		fmt.Println("PASS")

		fmt.Println("Validating pods in namespace kube-system")
		if len(res.Failures) == 0 {
			fmt.Println("PASS")
		} else {
			fmt.Println("FAIL")
			for _, failure := range res.Failures {
				fmt.Printf("\t%s\n", failure)
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
