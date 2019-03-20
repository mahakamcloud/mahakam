package network

import "net"

const (
	BridgeFormat = "clbr%s"
	TapDevFormat = "tap%sh%sh%s"
)

// ClusterHostGRE represents components at host to support
// cluster networking, i.e. GRE tunnel, bridge, tap device.
type ClusterHostGRE struct {
	BrName  string
	GREKey  string
	Tunnels []HostGRETunnel
}

type HostGRETunnel struct {
	TapDevName        string
	LocalHostIP       net.IP
	RemoteHostIP      net.IP
	TunnelReadyStatus bool
}
