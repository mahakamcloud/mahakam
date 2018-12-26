package provisioner

import (
	"fmt"

	"github.com/mahakamcloud/mahakam/pkg/node"
	log "github.com/sirupsen/logrus"
)

type CreateNode struct {
	Config node.NodeCreateConfig
	p      Provisioner
	log    log.FieldLogger
}

func NewCreateNode(config node.NodeCreateConfig, p Provisioner, log log.FieldLogger) *CreateNode {
	createNodeLog := log.WithField("task", fmt.Sprintf("create node in %s", config.Host))

	return &CreateNode{
		Config: config,
		p:      p,
		log:    createNodeLog,
	}
}

func (n *CreateNode) Run() error {
	err := n.p.CreateNode(n.Config)
	if err != nil {
		log.Errorf("error creating node '%v': %s", n.Config, err)
	}
	return nil
}
