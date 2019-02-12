package handlers

import (
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/network"
	"github.com/mahakamcloud/mahakam/pkg/provisioner"
	store "github.com/mahakamcloud/mahakam/pkg/resource_store"
	"github.com/sirupsen/logrus"
)

// Handlers holds common modules that each handler needs
type Handlers struct {
	AppConfig   *config.Config
	Store       store.ResourceStore
	Network     *network.NetworkManager
	Provisioner provisioner.Provisioner
	Log         logrus.FieldLogger
}

// New creates new handlers
func New(config *config.Config, provisioner provisioner.Provisioner, log logrus.FieldLogger) *Handlers {
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
		Log:         log,
	}
}
