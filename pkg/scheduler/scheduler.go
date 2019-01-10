package scheduler

import (
	"fmt"
	"net"

	"github.com/mahakamcloud/mahakam/pkg/config"
)

// Schedule interface defines a Scheduler that returns allocated node
type Schedule interface {
	GetHost(hosts []config.Host) (net.IP, error)
}

// GetHost return a single host IP
func GetHost(hosts []config.Host) (net.IP, error) {
	if len(hosts) == 0 {
		return nil, fmt.Errorf("Empty hosts config")
	}

	host := net.ParseIP(hosts[0].IPAddress)
	if host == nil {
		return nil, fmt.Errorf("Invalid host address")
	}

	return host, nil
}
