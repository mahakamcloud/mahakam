package handlers

import (
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/client"
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/network"
	"github.com/mahakamcloud/mahakam/pkg/provisioner"
	store "github.com/mahakamcloud/mahakam/pkg/resource_store"
	log "github.com/sirupsen/logrus"
)

// Handlers holds common modules that each handler needs
type Handlers struct {
	Store       store.ResourceStore
	Network     *network.NetworkManager
	Provisioner provisioner.Provisioner
}

// New creates new handlers
func New(storeConfig config.StorageBackendConfig, networkConfig config.NetworkConfig,
	provisioner provisioner.Provisioner) *Handlers {
	rs, err := store.New(storeConfig)
	if err != nil {
		log.Fatalf("Error initializing resource store in handlers: %s", err)
		return nil
	}

	n, err := network.New(rs, networkConfig)
	if err != nil {
		log.Fatalf("Error initializing network in handlers: %s", err)
		return nil
	}

	return &Handlers{
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
