package algorithm

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mahakamcloud/mahakam/pkg/config"
)

// RandomAllocator allocates a host random host from a list of available hosts
func RandomAllocator(hosts []config.Host) (config.Host, error) {
	rand.Seed(time.Now().UnixNano())

	if len(hosts) < 1 {
		return config.Host{}, fmt.Errorf("host list length is less than 1")
	}

	host := hosts[rand.Intn(len(hosts))]
	return host, nil
}
