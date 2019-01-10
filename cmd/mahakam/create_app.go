package main

import (
	"fmt"
	"os"

	"github.com/go-openapi/swag"

	"github.com/golang/glog"
	"github.com/mahakamcloud/mahakam/pkg/api/v1"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/client/apps"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/handlers"
	"github.com/spf13/cobra"
)

type CreateAppOptions struct {
	Name        string
	Owner       string
	ClusterName string
	ChartURL    string
	ChartValues string

	AppAPI v1.AppAPI
}

var cao = &CreateAppOptions{}

var createAppCmd = &cobra.Command{
	Use:   "app",
	Short: "Create application",
	Long:  `Create application on kubernetes cluster with one command`,
	Run: func(cmd *cobra.Command, args []string) {
		if cao.Name == "" {
			exitWithHelp(cmd, "Please provide name for your application.")
		}

		if cao.ClusterName == "" {
			exitWithHelp(cmd, "Please provide which cluster do you want to run the application.")
		}

		if cao.ChartURL == "" {
			exitWithHelp(cmd, "Please provide which helm chart to deploy your application.")
		}

		if cao.Owner == "" {
			// Hack since we don't have login mechanism yet
			cao.Owner = config.ResourceOwnerGojek
		}

		cao.AppAPI = handlers.GetMahakamAppClient(os.Getenv("MAHAKAM_API_SERVER_HOST"))

		res, err := RunCreateApp(cao)
		if err != nil {
			glog.Exit(err)
		}

		fmt.Println("Creating your application...")
		fmt.Printf("\nName:\t%s", swag.StringValue(res.Name))
		fmt.Printf("\nCluster:\t%s", res.ClusterName)
		fmt.Printf("\nApp endpoint:\t%s", res.ServiceFQDN)
		fmt.Printf("\n\nUse 'mahakam describe app %s' to monitor the state of your application", swag.StringValue(res.Name))
	},
}

func RunCreateApp(cao *CreateAppOptions) (*models.App, error) {
	req := &models.App{
		Name:        swag.String(cao.Name),
		Owner:       cao.Owner,
		ClusterName: cao.ClusterName,
		ChartURL:    cao.ChartURL,
		ChartValues: cao.ChartValues,
	}

	if cao.ChartValues != "" {
		f, err := os.Open(cao.ChartValues)
		if err != nil {
			return nil, fmt.Errorf("error reading app values path '%v': '%v'", cao, err.Error())
		}

		_, err = cao.AppAPI.UploadAppValues(apps.NewUploadAppValuesParams().
			WithValues(f).
			WithAppName(req.Name).
			WithOwner(swag.String(req.Owner)).
			WithClusterName(swag.String(req.ClusterName)))
		if err != nil {
			return nil, fmt.Errorf("error uploading app values '%v': '%v'", cao, err.Error())
		}
	}

	res, err := cao.AppAPI.CreateApp(apps.NewCreateAppParams().WithBody(req))
	if err != nil {
		return nil, fmt.Errorf("error creating app '%v': '%v'", cao, err.Error())
	}

	return res.Payload, nil
}

func init() {
	// Required flags
	createAppCmd.Flags().StringVarP(&cao.Name, "app-name", "a", "", "Name for your application")
	createAppCmd.Flags().StringVarP(&cao.ClusterName, "cluster-name", "c", "", "Name of your kubernetes cluster")
	createAppCmd.Flags().StringVarP(&cao.ChartURL, "chart", "u", "", "Helm chart url to run your application")

	// Optional flags
	createAppCmd.Flags().StringVarP(&cao.Owner, "owner", "o", "", "Name for your owner or app group")
	createAppCmd.Flags().StringVarP(&cao.ChartValues, "values", "v", "", "Helm values to override default one in the chart")

	createCmd.AddCommand(createAppCmd)
}
