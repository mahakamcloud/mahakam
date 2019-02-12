package client

import (
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/mahakamcloud/mahakam/pkg/api/v1"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/client"
	"github.com/mahakamcloud/mahakam/pkg/config"
)

func GetMahakamClient(host string) *client.Mahakam {
	t := httptransport.New(host, config.MahakamAPIBasePath, nil)
	c := client.New(t, strfmt.Default)
	return c
}

func GetMahakamClusterClient(host string) v1.ClusterAPI {
	t := httptransport.New(host, config.MahakamAPIBasePath, nil)
	c := client.New(t, strfmt.Default)
	return c.Clusters
}

func GetMahakamAppClient(host string) v1.AppAPI {
	t := httptransport.New(host, config.MahakamAPIBasePath, nil)
	c := client.New(t, strfmt.Default)
	return c.Apps
}
