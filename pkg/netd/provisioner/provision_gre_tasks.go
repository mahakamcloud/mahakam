package provisioner

import (
	"fmt"
	"net"

	"github.com/digitalocean/go-openvswitch/ovs"
	"github.com/sirupsen/logrus"
)

type CreateTapDev struct {
	tapDevName string
	log        logrus.FieldLogger
}

func NewCreateTapDev(tapDevName string, log logrus.FieldLogger) *CreateTapDev {
	tdLog := log.WithField("task", fmt.Sprintf("create tap dev %q", tapDevName))

	return &CreateTapDev{
		tapDevName: tapDevName,
		log:        tdLog,
	}
}

func (td *CreateTapDev) Run() error {
	td.log.Debugf("Creating tap device")

	// TODO(giri): implement create tap dev
	return nil
}

type CreateGREBridge struct {
	bridgeName string
	tapDevName string
	greKey     string
	remoteIP   net.IP
	ovsClient  *ovs.Client
	log        logrus.FieldLogger
}

func NewCreateGREBridge(bridgeName, tapDevName, greKey string, remoteIP net.IP, ovsClient *ovs.Client, log logrus.FieldLogger) *CreateGREBridge {
	greLog := log.WithField("task", fmt.Sprintf("create gre bridge %q with key %q and tap %q", bridgeName, greKey, tapDevName))

	return &CreateGREBridge{
		bridgeName: bridgeName,
		tapDevName: tapDevName,
		greKey:     greKey,
		remoteIP:   remoteIP,
		ovsClient:  ovsClient,
		log:        greLog,
	}
}

func (ob *CreateGREBridge) Run() error {
	if !ob.bridgeExists() {
		ob.log.Debugf("GRE bridge doesn't exist, creating one")
		if err := ob.ovsClient.VSwitch.AddBridge(ob.bridgeName); err != nil {
			return fmt.Errorf("failed to add new ovs bridge: %s", err)
		}
	}

	ob.log.Debugf("Attach tap dev %q to GRE bridge %q", ob.tapDevName, ob.bridgeName)
	if err := ob.ovsClient.VSwitch.AddPort(ob.bridgeName, ob.tapDevName); err != nil {
		return fmt.Errorf("failed to attach tap dev to ovs bridge: %s", err)
	}

	opts := ovs.InterfaceOptions{
		Type:     ovs.InterfaceTypeGRE,
		RemoteIP: ob.remoteIP.String(),
		Key:      ob.greKey,
	}
	ob.log.Debugf("Set GRE interface for bridge %q with key %q", ob.bridgeName, ob.greKey)
	if err := ob.ovsClient.VSwitch.Set.Interface(ob.tapDevName, opts); err != nil {
		return fmt.Errorf("failed to set gre interface: %s", err)
	}

	return nil
}

func (ob *CreateGREBridge) bridgeExists() bool {
	bridges, _ := ob.ovsClient.VSwitch.ListBridges()
	for _, br := range bridges {
		if ob.bridgeName == br {
			return true
		}
	}
	return false
}
