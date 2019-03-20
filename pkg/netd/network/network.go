package network

import (
	"net"
	"regexp"
)

const (
	// BridgeFormat represents bridge name formatted with
	// "clbr[GREKey]". For example, a bridge named clbr1 has GRE key 1.
	BridgeFormat = "clbr%s"
	// TapDevFormat represents tap dev name formatted with
	// "tap[GREKey]h[localIP]h[remoteIP]".
	TapDevFormat = "tap%sh%sh%s"
)

// ClusterHostGRE represents components at host to support
// cluster networking, i.e. GRE tunnel, bridge, tap device.
type ClusterHostGRE struct {
	ClusterName string
	BrName      string
	GREKey      string
	Tunnels     []HostGRETunnel
}

type HostGRETunnel struct {
	TapDevName        string
	LocalHostIP       net.IP
	RemoteHostIP      net.IP
	TunnelReadyStatus bool
}

func ParseGREKey(brname string) string {
	re := regexp.MustCompile("[0-9]+")
	strs := re.FindAllString(brname, -1)
	if len(strs) == 0 {
		return ""
	}
	return strs[len(strs)-1]
}
