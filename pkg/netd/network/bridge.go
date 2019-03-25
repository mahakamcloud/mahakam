package network

import (
	"github.com/digitalocean/go-openvswitch/ovs"
)

// NewBridge creates a new network bridge if it doesn't exists.
func createBridge(bridgeName string) error {
	o := ovs.New(
		ovs.Sudo(),
	)

	return o.VSwitch.AddBridge(bridgeName)
}
