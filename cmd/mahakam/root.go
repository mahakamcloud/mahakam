package main

import (
	"fmt"
	"os"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/client"
	"github.com/spf13/cobra"
)

const basePath = "/v1"

var RootCmd = &cobra.Command{
	Use:   "mahakam",
	Short: "application cloud platform",
	Long:  `Simple application cloud platform on Kubernetes`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func exitWithHelp(cmd *cobra.Command, err string) {
	fmt.Fprintln(os.Stderr, err)
	cmd.Help()
	os.Exit(1)
}

func GetClusterClient() *client.Mahakam {
	t := httptransport.New(os.Getenv("MAHAKAM_API_SERVER_HOST"), basePath, nil)
	c := client.New(t, strfmt.Default)
	return c
}
