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

	ipaddress, err := hostIP(nd.HostBridgeName)
	if err != nil {
		nd.Log.Errorf("error getting host ip address: %v", err)
		return
	}

	provisionAgent := agent.NewProvisionAgent(hostname, ipaddress.String(), nd.MahakamAPIServer, nd.Log)
	provisionAgent.Run()
}

func hostIP(brName string) (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		if iface.Name == brName {
			if addrs, _ := iface.Addrs(); len(addrs) > 0 {
				ip, _, _ := net.ParseCIDR(addrs[0].String())
				return ip, nil
			}
			return nil, fmt.Errorf("host bridge %q doesn't have IP", brName)
		}
	}
	return nil, fmt.Errorf("host bridge %q not found", brName)
}