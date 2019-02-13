package scheduler_test

import (
	"testing"

	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/scheduler"
	"github.com/stretchr/testify/assert"
)

func TestGetHostReturnsErrorForEmptyHostList(t *testing.T) {

	hosts := []config.Host{}
	_, err := scheduler.GetHost(hosts)

	assert.NotNil(t, err)
	assert.Equal(t, "empty hosts config", err.Error(), "they should be equal")
}

func TestGetHostReturnsErrorForInvalidHost(t *testing.T) {

	hosts := []config.Host{}
	_, err := scheduler.GetHost(hosts)

	assert.NotNil(t, err)
	assert.Equal(t, "empty hosts config", err.Error(), "they should be equal")
}
