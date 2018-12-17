package handlers

import (
	"github.com/mahakamcloud/mahakam/pkg/config"
	store "github.com/mahakamcloud/mahakam/pkg/resource_store"
	log "github.com/sirupsen/logrus"
)

// Handlers holds common modules that each handler needs
type Handlers struct {
	Store store.ResourceStore
}

// New creates new handlers
func New(storeConfig config.StorageBackendConfig) *Handlers {
	rs, err := store.New(storeConfig)
	if err != nil {
		log.Fatalf("Error initializing handlers: %s", err)
		return nil
	}

	return &Handlers{
		Store: rs,
	}
}
