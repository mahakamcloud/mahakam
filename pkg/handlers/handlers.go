package handlers

import (
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/mahakamcloud/mahakam/pkg/api/v1"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/client"
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/network"
	"github.com/mahakamcloud/mahakam/pkg/provisioner"
	store "github.com/mahakamcloud/mahakam/pkg/resource_store"
	log "github.com/sirupsen/logrus"
)

// Handlers holds common modules that each handler needs
type Handlers struct {
	AppConfig   *config.Config
	Store       store.ResourceStore
	Network     *network.NetworkManager
	Provisioner provisioner.Provisioner
}

// New creates new handlers
func New(config *config.Config, provisioner provisioner.Provisioner) *Handlers {
	rs, err := store.New(config.KVStoreConfig)
	if err != nil {
		log.Fatalf("Error initializing resource store in handlers: %s", err)
		return nil
	}

	n, err := network.New(rs, config.NetworkConfig)
	if err != nil {
		log.Fatalf("Error initializing network in handlers: %s", err)
		return nil
	}

	return &Handlers{
		AppConfig:   config,
		Store:       rs,
		Network:     n,
		Provisioner: provisioner,
	}
}

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
