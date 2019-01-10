package scheduler_test

import (
	"net"
	"testing"

	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/scheduler"
	"github.com/stretchr/testify/assert"
)

func TestGetHostReturnsFirstHost(t *testing.T) {
	host1 := config.Host{Name: "i-test-01", IPAddress: "127.0.0.1"}
	host2 := config.Host{Name: "i-test-02", IPAddress: "127.0.1.1"}

	hostConfig := []config.Host{
		host1,
		host2,
	}

	scheduledHost, err := scheduler.GetHost(hostConfig)

	expectedHost := net.ParseIP("127.0.0.1")

	assert.Nil(t, err)
	assert.Equal(t, scheduledHost, expectedHost, "they should be equal")
}

func TestGetHostReturnsErrorForEmptyHostList(t *testing.T) {

	hosts := []config.Host{}
	_, err := scheduler.GetHost(hosts)

	assert.NotNil(t, err)
	assert.Equal(t, "Empty hosts config", err.Error(), "they should be equal")
}

func TestGetHostReturnsErrorForInvalidHost(t *testing.T) {

	hosts := []config.Host{}
	_, err := scheduler.GetHost(hosts)

	assert.NotNil(t, err)
	assert.Equal(t, "Empty hosts config", err.Error(), "they should be equal")
}
