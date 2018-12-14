package resourcestore

import (
	"fmt"
	"regexp"
	"time"
)

// Resource is base interface for all stored resources or objects
type Resource interface {
	GetResource() *BaseResource
	PreCheck() error
}

// BaseResource is the base struct for all stored resources or objects
type BaseResource struct {
	Name         string    `json:"name"`
	Kind         string    `json:"kind"`
	Owner        string    `json:"owner"`
	CreatedTime  time.Time `json:"createdTime,omitempty"`
	ModifiedTime time.Time `json:"modifiedTime,omitempty"`
	Revision     uint64    `json:"revision"`
	Status       Status    `json:"status"`
}

func (br *BaseResource) GetResource() *BaseResource {
	return br
}

func (br *BaseResource) PreCheck() error {
	if br.Owner == "" {
		return fmt.Errorf("BaseResource owner attribute cannot be empty")
	}

	var validName = regexp.MustCompile(`^[\w\d\-]+$`)
	if !validName.MatchString(br.Name) {
		return fmt.Errorf("BaseResource name %s is invalid", br.Name)
	}
	return nil
}

func (br *BaseResource) BuildKey() string {
	return fmt.Sprintf("%s/%s/%s", br.Kind, br.Owner, br.Name)
}
