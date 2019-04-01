package resourcestore

import (
	"encoding/json"
	"fmt"

	"github.com/docker/libkv/store"
	"github.com/mahakamcloud/mahakam/pkg/model"
	"github.com/mahakamcloud/mahakam/pkg/resource_store/resource"
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

// Add adds new resource to kv store
func (kvr *kvResourceStore) AddV1(r model.ResourceBuilder) (id string, err error) {
	if err := r.Validate(); err != nil {
		return "", fmt.Errorf("kv resource precheck failed: %s", err)
	}

	// TODO(giri): check if key exists or duplicated
	key := r.BuildKey()
	r.AddMetadata()

	value, err := r.MarshalJSON()
	if err != nil {
		return "", fmt.Errorf("add kv resource serialization error: %s", err)
	}

	opts := &store.WriteOptions{
		IsDir: false,
	}
	_, _, err = kvr.store.AtomicPut(key, value, nil, opts)
	if err != nil {
		return "", fmt.Errorf("add kv resource atomic put error: %s", err)
	}

	// TODO(giri): Implement optimistic lock with revision.

	return r.GetID(), nil
}

// Get retrieves single resource with values from kv store,
// must include owner, name, and kind in the passed resource struct
func (kvr *kvResourceStore) GetV1(r model.ResourceBuilder) error {
	res, err := kvr.store.Get(r.BuildKey())
	if err != nil {
		return fmt.Errorf("error getting response from kv store: %s", err)
	}

	err = r.UnmarshalJSON(res.Value)
	if err != nil {
		return fmt.Errorf("error unmarshalling resource: %s", err)
	}

	return nil
}

// List returns list of resource from store
func (kvr *kvResourceStore) ListV1(owner string, kind model.ResourceKind, list model.ResourceBuilderList) error {
	var resources []model.ResourceBuilder

	kvpairs, err := kvr.store.List(string(kind) + "/" + owner + "/")
	if err != nil && err != store.ErrKeyNotFound {
		return fmt.Errorf("error getting list of kvpairs from path: %s", err)
	}

	for _, kv := range kvpairs {
		r := list.ResourceBuilder()

		err = r.UnmarshalJSON(kv.Value)
		if err != nil {
			return fmt.Errorf("error unmarshalling resources: %s", err)
		}

		resources = append(resources, r)
	}

	list.WithItems(resources)
	return nil
}

// Add adds new resource to kv store
func (kvr *kvResourceStore) Add(r resource.Resource) (id string, err error) {
	if err := r.PreCheck(); err != nil {
		return "", fmt.Errorf("kv resource precheck failed: %s", err)
	}

	// TODO(giri): check if key exists or duplicated
	key := r.BuildKey()
	r.BuildResource()

	value, err := json.Marshal(r)
	if err != nil {
		return "", fmt.Errorf("add kv resource serialization error: %s", err)
	}

	opts := &store.WriteOptions{
		IsDir: false,
	}
	_, res, err := kvr.store.AtomicPut(key, value, nil, opts)
	if err != nil {
		return "", fmt.Errorf("add kv resource atomic put error: %s", err)
	}

	r.GetResource().Revision = res.LastIndex
	return r.GetResource().ID, nil
}

// Get retrieves single resource with values from kv store,
// must include owner, name, and kind in the passed resource struct
func (kvr *kvResourceStore) Get(r resource.Resource) error {
	if r.GetResource().Owner == "" {
		return fmt.Errorf("owner parameter is required for getting resource")
	}
	if r.GetResource().Name == "" {
		return fmt.Errorf("name parameter is required for getting resource")
	}
	if r.GetResource().Kind == "" {
		return fmt.Errorf("kind parameter is required for getting resource")
	}

	res, err := kvr.store.Get(r.BuildKey())
	if err != nil {
		return fmt.Errorf("error getting response from kv store: %s", err)
	}

	err = json.Unmarshal(res.Value, r)
	if err != nil {
		return fmt.Errorf("error unmarshalling resource: %s", err)
	}

	// TODO(giri): filter based on given labels and scope
	r.GetResource().Revision = res.LastIndex
	return nil
}

func (kvr *kvResourceStore) List(owner string, kind resource.ResourceKind, list resource.ResourceList) error {
	var resources []resource.Resource

	kvpairs, err := kvr.store.List(string(kind) + "/" + owner + "/")
	if err != nil && err != store.ErrKeyNotFound {
		return fmt.Errorf("error getting list of kvpairs from path: %s", err)
	}

	for _, kv := range kvpairs {
		r := list.Resource()

		err = json.Unmarshal(kv.Value, r)
		if err != nil {
			return fmt.Errorf("error unmarshalling resources: %s", err)
		}

		resources = append(resources, r)
	}

	list.WithItems(resources)
	return nil
}

func (kvr *kvResourceStore) Update(resource resource.Resource) (revision int64, err error) {
	if err := resource.PreCheck(); err != nil {
		return -1, fmt.Errorf("kv resource precheck failed: %s", err)
	}
	resource.UpdateResource()

	value, err := json.Marshal(resource)
	if err != nil {
		return -1, fmt.Errorf("add kv resource serialization error: %s", err)
	}

	prev := &store.KVPair{
		Key:       resource.BuildKey(),
		LastIndex: resource.GetResource().Revision,
	}
	opts := &store.WriteOptions{
		IsDir: false,
	}
	_, res, err := kvr.store.AtomicPut(resource.BuildKey(), value, prev, opts)
	if err != nil {
		return -1, fmt.Errorf("add kv resource atomic put error: %s", err)
	}

	resource.GetResource().Revision = res.LastIndex
	return int64(res.LastIndex), nil
}

func (kvr *kvResourceStore) Delete(owner string, id string, resource resource.Resource) error {
	fmt.Println("libkvResourceStore Delete method not implemented")
	return nil
}

// AddFromPath adds new resource to kv store with specific key path
func (kvr *kvResourceStore) AddFromPath(path string, r resource.Resource) (id string, err error) {
	if path == "" {
		return "", fmt.Errorf("must provide non-empty path for storing resource: %s", err)
	}
	if err := r.PreCheck(); err != nil {
		return "", fmt.Errorf("kv resource precheck failed: %s", err)
	}
	r.BuildResource()

	value, err := json.Marshal(r)
	if err != nil {
		return "", fmt.Errorf("add kv resource serialization error: %s", err)
	}

	opts := &store.WriteOptions{
		IsDir: false,
	}
	_, res, err := kvr.store.AtomicPut(path, value, nil, opts)
	if err != nil {
		return "", fmt.Errorf("add kv resource atomic put error: %s", err)
	}

	r.GetResource().Revision = res.LastIndex
	return r.GetResource().ID, nil
}

// Get retrieves single resource from specified key path
func (kvr *kvResourceStore) GetFromPath(path string, r resource.Resource) error {
	if path == "" {
		return fmt.Errorf("must provide non-empty path for getting resource")
	}

	res, err := kvr.store.Get(path)
	if err != nil {
		return fmt.Errorf("error getting response from kv store: %s", err)
	}

	err = json.Unmarshal(res.Value, r)
	if err != nil {
		return fmt.Errorf("error unmarshalling resource: %s", err)
	}

	r.GetResource().Revision = res.LastIndex
	return nil
}

// ListFromPath is quick hack to retrieve list keys from given path
func (kvr *kvResourceStore) ListFromPath(path string, filter Filter, resources resource.ResourceList) error {
	kvpairs, err := kvr.store.List(path)
	if err != nil && err != store.ErrKeyNotFound {
		return fmt.Errorf("error getting list of kvpairs from path: %s", err)
	}

	var items []resource.Resource
	for _, kvpair := range kvpairs {
		r := resources.Resource()
		err := json.Unmarshal(kvpair.Value, r)
		if err != nil {
			return fmt.Errorf("error unmarshalling resource: %s", err)
		}

		if filter != nil {
			ok := ApplyFilter(filter, r)
			if err != nil {
				return fmt.Errorf("error filtering resources: %s", err)
			}
			if !ok {
				continue
			}
		}

		items = append(items, r)
	}
	resources.WithItems(items)
	return nil
}

// ListKeysFromPath is quick hack to retrieve list keys from given path
func (kvr *kvResourceStore) ListKeysFromPath(path string) ([]string, error) {
	var keys []string

	kvpairs, err := kvr.store.List(path)
	if err != nil && err != store.ErrKeyNotFound {
		return []string{}, fmt.Errorf("error getting list of kvpairs from path: %s", err)
	}

	for _, kvpair := range kvpairs {
		keys = append(keys, kvpair.Key)
	}
	return keys, nil
}

func (kvr *kvResourceStore) UpdateFromPath(path string, r resource.Resource) (revision int64, err error) {
	if path == "" {
		return -1, fmt.Errorf("must provide non-empty path for updating resource: %s", err)
	}
	if err := r.PreCheck(); err != nil {
		return -1, fmt.Errorf("kv resource precheck failed: %s", err)
	}
	r.UpdateResource()

	value, err := json.Marshal(r)
	if err != nil {
		return -1, fmt.Errorf("add kv resource serialization error: %s", err)
	}

	prev := &store.KVPair{
		Key:       path,
		LastIndex: r.GetResource().Revision,
	}
	opts := &store.WriteOptions{
		IsDir: false,
	}
	_, res, err := kvr.store.AtomicPut(path, value, prev, opts)
	if err != nil {
		return -1, fmt.Errorf("add kv resource atomic put error: %s", err)
	}

	r.GetResource().Revision = res.LastIndex
	return int64(res.LastIndex), nil
}

func (kvr *kvResourceStore) DeleteFromPath(path string) error {
	fmt.Println("libkvResourceStore DeleteFromPath method not implemented")
	return nil
}

func (kvr *kvResourceStore) KeyExists(path string) bool {
	exists, err := kvr.store.Exists(path)
	if err != nil && err != store.ErrKeyNotFound {
		return true
	}
	return exists
}
