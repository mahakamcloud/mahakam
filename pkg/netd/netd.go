package netd

import (
	"fmt"
	"net"
	"os"

	"github.com/sirupsen/logrus"
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
	_, err := os.Hostname()
	if err != nil {
		nd.Log.Errorf("error getting hostname: %v", err)
		return
	}

	_, err = getBridgeIpAddr(nd.HostBridgeName)
	if err != nil {
		nd.Log.Errorf("error getting host ip address: %v", err)
		return
	}
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
