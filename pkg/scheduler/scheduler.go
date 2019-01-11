package scheduler

import (
	"fmt"
	"net"

	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/scheduler/algorithm"
)

// Scheduler interface defines a Scheduler that returns allocated node
type Scheduler interface {
	GetHost(hosts []config.Host) (net.IP, error)
}

// GetHost return a single host IP
func GetHost(hosts []config.Host) (net.IP, error) {
	if len(hosts) == 0 {
		return nil, fmt.Errorf("Empty hosts config")
	}

	host, err := algorithm.RandomAllocator(hosts)
	if err != nil {
		return nil, fmt.Errorf("Empty hosts config")
	}

	hostIP := net.ParseIP(host.IPAddress)
	if hostIP == nil {
		return nil, fmt.Errorf("Invalid host address")
	}

	return hostIP, nil
}
