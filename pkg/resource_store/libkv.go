package resourcestore

import (
	"fmt"

	"github.com/docker/libkv/store"
)

type libkvResourceStore struct {
	kv store.Store
}

func newLibkv(kv store.Store) ResourceStore {
	return &libkvResourceStore{
		kv: kv,
	}
}

func (r *libkvResourceStore) Add(resource Resource) (id string, err error) {
	fmt.Println("libkvResourceStore Add method not implemented")
	return "", nil
}

func (r *libkvResourceStore) Get(owner string, key string, resource Resource) error {
	fmt.Println("libkvResourceStore Get method not implemented")
	return nil
}

func (r *libkvResourceStore) List(owner string, resources interface{}) error {
	fmt.Println("libkvResourceStore List method not implemented")
	return nil
}

func (r *libkvResourceStore) Update(resource Resource) (revision int64, err error) {
	fmt.Println("libkvResourceStore Update method not implemented")
	return 0, nil
}

func (r *libkvResourceStore) Delete(owner string, id string, resource Resource) error {
	fmt.Println("libkvResourceStore Delete method not implemented")
	return nil
}
