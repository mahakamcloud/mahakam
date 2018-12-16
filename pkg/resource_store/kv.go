package resourcestore

import (
	"encoding/json"
	"fmt"

	"github.com/docker/libkv/store"
)

type kvResourceStore struct {
	store store.Store
}

// NewKVResourceStore creates key value store for resources
func NewKVResourceStore(s store.Store) ResourceStore {
	return &kvResourceStore{
		store: s,
	}
}

func (kvr *kvResourceStore) Add(r Resource) (id string, err error) {
	if err := r.PreCheck(); err != nil {
		return "", fmt.Errorf("KV resource precheck failed: %s", err)
	}

	// TODO(giri): check if key exists or duplicated
	key := r.BuildKey()

	value, err := json.Marshal(r.BuildResource())
	if err != nil {
		return "", fmt.Errorf("Add KV resource serialization error: %s", err)
	}

	opts := &store.WriteOptions{
		IsDir: false,
	}
	_, res, err := kvr.store.AtomicPut(key, value, nil, opts)
	if err != nil {
		return "", fmt.Errorf("Add KV resource atomic put error: %s", err)
	}

	r.GetResource().Revision = res.LastIndex
	return r.GetResource().ID, nil
}

func (kvr *kvResourceStore) Get(owner string, name string, resource Resource) error {
	err := kvr.find(owner, name, nil)
	if err != nil {
		return fmt.Errorf("Error getting resource from kv store: %s", err)
	}

	return nil
}

func (kvr *kvResourceStore) find(owner string, name string, resource Resource) error {
	if owner == "" {
		return fmt.Errorf("Owner parameter is required for finding resource")
	}
	if name == "" {
		return fmt.Errorf("Name parameter is required for finding resource")
	}

	res, err := kvr.store.Get(resource.BuildKey())
	if err != nil {
		return fmt.Errorf("Error getting resource from kv store: %s", err)
	}

	err = json.Unmarshal(res.Value, resource)
	if err != nil {
		return fmt.Errorf("Error unmarshalling resource: %s", err)
	}

	// TODO(giri): filter based on given labels and scope
	resource.GetResource().Revision = res.LastIndex
	return nil
}

func (kvr *kvResourceStore) List(owner string, resources interface{}) error {
	fmt.Println("libkvResourceStore List method not implemented")
	return nil
}

func (kvr *kvResourceStore) Update(resource Resource) (revision int64, err error) {
	fmt.Println("libkvResourceStore Update method not implemented")
	return 0, nil
}

func (kvr *kvResourceStore) Delete(owner string, id string, resource Resource) error {
	fmt.Println("libkvResourceStore Delete method not implemented")
	return nil
}
