package provisioner_test

import (
	"fmt"
	"io/ioutil"
	"net"
	"testing"

	"github.com/mahakamcloud/mahakam/pkg/network"

	"github.com/mahakamcloud/mahakam/pkg/utils"

	"github.com/stretchr/testify/assert"

	"github.com/sirupsen/logrus"

	"github.com/mahakamcloud/mahakam/pkg/node"
	. "github.com/mahakamcloud/mahakam/pkg/provisioner"

	gomock "github.com/golang/mock/gomock"
)

var (
	n = node.NodeCreateConfig{
		Host: net.ParseIP("10.10.10.10"),
	}

	cn = &network.ClusterNetwork{
		Name: "fake cluster network",
	}
)

func nilLogger() *logrus.Logger {
	l := logrus.New()
	l.Out = ioutil.Discard

	return l
}

func TestCreateNode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	l := nilLogger()
	p := NewMockProvisioner(ctrl)

	tests := []struct {
		name        string
		expectError error
	}{
		{
			name:        "test create node where provisioner successfully runs",
			expectError: nil,
		},
		{
			name:        "test create node where provisioner fails",
			expectError: fmt.Errorf("fake error from provisioner"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p.EXPECT().CreateNode(gomock.Any()).Return(test.expectError)

			cn := NewCreateNode(n, p, l)
			err := cn.Run()

			assert.Equal(t, test.expectError, err)
		})
	}
}

func TestCheckClusterNetworkNodes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		expectError error
		expectReady bool
	}{
		{
			name:        "test check cluster network nodes where gateway is ready",
			expectError: nil,
			expectReady: true,
		},
		{
			name:        "test check cluster network nodes where gateway is not ready",
			expectError: fmt.Errorf("timeout waiting for cluster gateway to be up"),
			expectReady: false,
		},
	}

	l := nilLogger()
	pc := utils.NewMockPingChecker(ctrl)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pc.EXPECT().ICMPPingNWithDelay(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(test.expectReady)

			c := NewCheckClusterNetworkNodes(cn, l, pc)
			err := c.Run()

			if err != nil {
				assert.Error(t, err, test.expectError)
			}

			if err == nil {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckNode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		nodeIP      net.IP
		expectError error
		expectReady bool
	}{
		{
			name:        "test check node where node is ready",
			nodeIP:      net.ParseIP("1.2.3.4"),
			expectError: nil,
			expectReady: true,
		},
		{
			name:        "test check node where node is not ready",
			nodeIP:      net.ParseIP("1.2.3.4"),
			expectError: fmt.Errorf("timeout waiting for node to be up"),
			expectReady: false,
		},
	}

	l := nilLogger()
	pc := utils.NewMockPingChecker(ctrl)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pc.EXPECT().ICMPPingNWithDelay(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(test.expectReady)

			c := NewCheckNode(test.nodeIP, l, pc)
			err := c.Run()

			if err != nil {
				assert.Error(t, err, test.expectError)
			}

			if err == nil {
				assert.NoError(t, err)
			}
		})
	}
}
