package resourcestore

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx/types"
	uuid "github.com/satori/go.uuid"
)

// PostgresResource is the base struct for all stored postgres row
type PostgresResource struct {
	Key          string         `db:"key"`
	ID           string         `db:"id"`
	Name         string         `db:"name"`
	Kind         string         `db:"kind"`
	Owner        string         `db:"owner"`
	CreatedTime  time.Time      `db:"created_time"`
	ModifiedTime time.Time      `db:"modified_time"`
	Revision     uint64         `db:"revision"`
	Status       Status         `db:"status"`
	Value        types.JSONText `db:"value"`
}

func NewPostgresResource(r Resource) (*PostgresResource, error) {
	br := r.GetResource()
	now := time.Now()
	id := uuid.NewV4().String()

	pr := &PostgresResource{
		Key:          br.BuildKey(),
		ID:           id,
		Name:         br.Name,
		Kind:         br.Kind,
		Owner:        br.Owner,
		CreatedTime:  now,
		ModifiedTime: now,
		Revision:     br.Revision,
		Status:       br.Status,
	}

	v, err := json.Marshal(br)
	if err != nil {
		return nil, fmt.Errorf("PostgresResource error marshaling base resource")
	}

	err = pr.Value.Scan(v)
	if err != nil {
		return nil, fmt.Errorf("PostgresResource error scanning base resource")
	}

	return pr, nil
}
