package scheduler

import (
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

	host := net.ParseIP(hosts[0].IPAddress)
	return host, nil
}
