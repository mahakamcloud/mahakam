package resourcestore

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mahakamcloud/mahakam/pkg/config"
	// sqlx requires postgres driver registration with this lib
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	SQLCreateTable = `
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
		value 					JSONB
	)
	`

	SQLInsertResource = `
	INSERT INTO resource
		(key, id, name, kind, owner, created_time, modified_time, revision, status, value)
	VALUES
		(:key, :id, :name, :kind, :owner, :created_time, :modified_time, :revision, :status, :value)
	`
)

type postgresResourceStore struct {
	db *sqlx.DB
}

// NewPostgresResourceStore creates new resource store with postgres backend
func NewPostgresResourceStore(c config.StorageBackendConfig) (ResourceStore, error) {
	conn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", c.Username,
		c.Password, c.Address, c.Bucket)
	log.Debugf("Creating postgres client with connection string: %s", conn)

	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		log.Errorf("Error connecting to postgres: %s", err)
		return nil, fmt.Errorf("Error connecting to postgres: %s", err)
	}
	p := &postgresResourceStore{
		db: db,
	}

	if _, err := p.db.Exec(SQLCreateTable); err != nil {
		return nil, fmt.Errorf("Postgres failed to create table: %s", err)
	}

	return p, nil
}

func (p *postgresResourceStore) Add(resource Resource) (id string, err error) {
	if err := resource.PreCheck(); err != nil {
		return "", fmt.Errorf("Postgres resource precheck failed: %s", err)
	}

	pr, err := NewPostgresResource(resource)
	if err != nil {
		return "", fmt.Errorf("Postgres resource creation error: %s", err)
	}

	_, err = p.db.NamedExec(SQLInsertResource, pr)
	if err != nil {
		return "", fmt.Errorf("Postgres resource insertion error: %s", err)
	}

	return pr.ID, nil
}

func (p *postgresResourceStore) Get(resource Resource) error {
	fmt.Println("postgresResourceStore Get method not implemented")
	return nil
}

func (p *postgresResourceStore) List(owner string, resources interface{}) error {
	fmt.Println("postgresResourceStore List method not implemented")
	return nil
}

func (p *postgresResourceStore) Update(resource Resource) (revision int64, err error) {
	fmt.Println("postgresResourceStore Update method not implemented")
	return 0, nil
}

func (p *postgresResourceStore) Delete(owner string, id string, resource Resource) error {
	fmt.Println("postgresResourceStore Delete method not implemented")
	return nil
}
