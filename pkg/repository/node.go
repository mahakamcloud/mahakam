package repository

import (
	"encoding/json"
	"fmt"

	"github.com/go-openapi/swag"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/kvstore"
)

const (
	// RoleNodeLabelValue represents role value of bare metal node
	RoleNodeLabelValue = "node"
)

type NodeRepository struct {
	store *kvstore.KVStore
}

func NewNodeRepository() (*NodeRepository, error) {
	store, err := kvstore.New()
	if err != nil {
		return nil, err
	}
	return &NodeRepository{store}, nil
}

func (r *NodeRepository) Put(b *models.Node) error {
	key := fmt.Sprintf("%s/%s/%s", b.Kind, b.Owner, swag.StringValue(b.Name))
	val, err := json.Marshal(b)
	if err != nil {
		return err
	}
	return r.store.Put(key, val)
}

func (r *NodeRepository) List() ([]*models.Node, error) {
	// TODO : Remove hard coded values
	key := fmt.Sprintf("%s/%s", "node", "mahakam")

	vals, err := r.store.List(key)
	if err != nil {
		return nil, err
	}

	nodes := make([]*models.Node, 0)

	for _, v := range vals {
		bm := &models.Node{}
		err := json.Unmarshal(v, bm)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling json value. Value: %v, Error: %v", string(v), err.Error())
		}
		nodes = append(nodes, bm)
	}
	return nodes, nil
}
