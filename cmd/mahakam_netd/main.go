package main

import (
	"fmt"
	"os"

	"github.com/mahakamcloud/mahakam/pkg/netd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var nd = &netd.NetDaemon{}

var app = &cobra.Command{
	Use:   "mahakam-netd",
	Short: "network daemon for mahakam cloud platform",
	Long:  `Network daemon for mahakam cloud platform`,
	Run: func(cmd *cobra.Command, args []string) {
		if nd.MahakamAPIServer == "" {
			exitWithHelp(cmd, "Please provide Mahakam API Server endpoint.")
		}

		logrus.Info("Mahakam netd starting...")
		netd.Run(nd)
	},
}

func main() {
	log := logrus.New().WithField("app", "mahakam-netd")
	log.Println("inititializing...")
	defer log.Println("exiting...")

	nd.Log = log

	registerFlags()

	if err := app.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func registerFlags() {
	// Required flags
	app.Flags().StringVarP(&nd.MahakamAPIServer, "mahakam-api-server", "a", "", "Endpoint of Mahakam API Server")

	// Optional flags
	app.Flags().StringVarP(&nd.HostBridgeName, "host-bridge", "b", "mahabr0", "Host bridge name that allows host to join Mahakam network")
}

func exitWithHelp(cmd *cobra.Command, err string) {
	fmt.Fprintln(os.Stderr, err)
	cmd.Help()
	os.Exit(1)
}
