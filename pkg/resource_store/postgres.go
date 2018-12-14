package resourcestore

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	// sqlx requires postgres driver registration with this lib
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type postgresResourceStore struct {
	db *sqlx.DB
}

// NewPostgresResourceStore creates new resource store with postgres backend
func NewPostgresResourceStore(config StorageBackendConfig) (ResourceStore, error) {
	conn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", config.Username,
		config.Password, config.Address, config.Bucket)
	log.Debugf("Creating postgres client with connection string: %s", conn)

	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		log.Errorf("Error connecting to postgres: %s", err)
		return nil, fmt.Errorf("Error connecting to postgres: %s", err)
	}
	p := &postgresResourceStore{
		db: db,
	}

	if err = p.createTable(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *postgresResourceStore) Add(ctx context.Context, resource Resource) (id string, err error) {
	fmt.Println("postgresResourceStore Add method not implemented")
	return "", nil
}

func (p *postgresResourceStore) Get(ctx context.Context, owner string, key string, resource Resource) error {
	fmt.Println("postgresResourceStore Get method not implemented")
	return nil
}

func (p *postgresResourceStore) List(ctx context.Context, owner string, resources interface{}) error {
	fmt.Println("postgresResourceStore List method not implemented")
	return nil
}

func (p *postgresResourceStore) Update(ctx context.Context, resource Resource) (revision int64, err error) {
	fmt.Println("postgresResourceStore Update method not implemented")
	return 0, nil
}

func (p *postgresResourceStore) Delete(ctx context.Context, owner string, id string, resource Resource) error {
	fmt.Println("postgresResourceStore Delete method not implemented")
	return nil
}

func (p *postgresResourceStore) createTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS resource (
		key 						TEXT PRIMARY KEY,
		id 							TEXT,
		name 						TEXT,
		kind						TEXT,
		owner 					TEXT,
		created_time 		TIMESTAMP,
		modified_time 	TIMESTAMP,
		revision 				BIGINT,
		status 					TEXT,
		labels					JSONB,
		value 					JSONB
	)
	`

	if _, err := p.db.Exec(sql); err != nil {
		return fmt.Errorf("Postgres failed to create table: %s", err)
	}

	return nil
}
