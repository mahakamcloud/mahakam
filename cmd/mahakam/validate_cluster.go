package main

import (
	"fmt"
	"os"

	"github.com/golang/glog"

	"github.com/mahakamcloud/mahakam/pkg/api/v1"
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

		err := RunValidateCluster(vco)
		if err != nil {
			glog.Exit(err)
		}

		// TODO(giri): print out validation result
		fmt.Println("Validating kubernetes cluster...")
	},
}

func RunValidateCluster(vco *ValidateClusterOptions) error {
	// TODO(giri): call validate cluster service
	return nil
}

func init() {
	// Required flags
	validateClusterCmd.Flags().StringVarP(&cco.Name, "cluster-name", "c", "", "Name for your kubernetes cluster")

	// Optional flags
	validateClusterCmd.Flags().StringVarP(&cco.Owner, "owner", "o", "", "Owner of your kubernetes cluster")

	createCmd.AddCommand(validateClusterCmd)
}
