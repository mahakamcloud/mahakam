package scheduler

import (
	"errors"
	"net"

	"github.com/mahakamcloud/mahakam/pkg/config"
)

// Schedule interface defines a Scheduler that returns allocated node
type Schedule interface {
	GetHost(hosts []config.HostsConfig) (net.IP, error)
}

// GetHost return a single host IP
func GetHost(hostConfig config.HostsConfig) (net.IP, error) {
	hosts := hostConfig.Hosts
	if len(hosts) == 0 {
		return nil, errors.New("Empty hosts config")
	}

	host := net.ParseIP(hosts[0].IPAddress)
	if host == nil {
		return nil, errors.New("Invalid host address")
	}

	return host, nil
}
