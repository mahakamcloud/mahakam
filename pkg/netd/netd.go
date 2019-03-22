package netd

import (
	"fmt"
	"net"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/mahakamcloud/mahakam/pkg/netd/agent"
)

type NetDaemon struct {
	MahakamAPIServer string
	HostBridgeName   string

	hostname  string
	ipaddress net.IP

	Log logrus.FieldLogger
}

func Run(nd *NetDaemon) {
	// TODO(giri): Self registration of host
	hostname, err := os.Hostname()
	if err != nil {
		nd.Log.Errorf("error getting hostname: %v", err)
		return
	}

	ipaddress, err := getBridgeIpAddr(nd.HostBridgeName)
	if err != nil {
		nd.Log.Errorf("error getting host ip address: %v", err)
		return
	}

	provisionAgent := agent.NewProvisionAgent(hostname, ipaddress.String(), nd.MahakamAPIServer, nd.Log)
	provisionAgent.Run()
}

func getBridgeIpAddr(bridgeName string) (net.IP, error) {
	iface, err := net.InterfaceByName(bridgeName)
	if err != nil {
		return nil, err
	}

	if addrs, _ := iface.Addrs(); len(addrs) > 0 {
		ip, _, _ := net.ParseCIDR(addrs[0].String())
		return ip, nil
	}
	return nil, fmt.Errorf("host bridge %q doesn't have IP", bridgeName)
}
