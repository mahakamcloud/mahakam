package network

// ClusterHostGRE represents components at host to support
// cluster networking, i.e. GRE tunnel, bridge, tap device.
type ClusterHostGRE struct {
	BrName            string
	TapDevName        string
	TunnelReadyStatus string
}
