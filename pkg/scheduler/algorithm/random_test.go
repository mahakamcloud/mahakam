package algorithm_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/scheduler/algorithm"
)

func TestRandomAllocatorReturnsEmptyHostForEmptyHostList(t *testing.T) {
	const hostLenErrorMsg = "host list length is less than 1"
	hostsConfig := []config.Host{}

	host, err := algorithm.RandomAllocator(hostsConfig)

	assert.Equal(t, err.Error(), hostLenErrorMsg, "they should be equal")
	assert.Equal(t, host, config.Host{}, "they should be equal")
}

func TestRandomAllocatorReturnsValidHostForNonEmptyHostList(t *testing.T) {
	host1 := config.Host{Name: "i-test-01", IPAddress: "127.0.0.1"}
	hostsConfig := []config.Host{host1}

	host, err := algorithm.RandomAllocator(hostsConfig)

	assert.Nil(t, err)
	assert.Equal(t, host, host1, "they should be equal")
}
