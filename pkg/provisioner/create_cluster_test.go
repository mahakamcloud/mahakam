package provisioner_test

import (
	"fmt"
	"io/ioutil"
	"net"
	"testing"

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
