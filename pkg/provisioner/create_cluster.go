package provisioner

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/network"
	"github.com/mahakamcloud/mahakam/pkg/node"
	store "github.com/mahakamcloud/mahakam/pkg/resource_store"
	"github.com/mahakamcloud/mahakam/pkg/resource_store/resource"
	"github.com/mahakamcloud/mahakam/pkg/utils"
	"github.com/mahakamcloud/mahakam/pkg/validation"
	"github.com/sirupsen/logrus"
)

type PreCreateCheck struct {
	clustername    string
	clusterKeyPath string
	log            logrus.FieldLogger
	store          store.ResourceStore
}

func NewPreCreateCheck(clustername string, log logrus.FieldLogger, store store.ResourceStore) *PreCreateCheck {
	preCreateCheckLog := log.WithField("cluster", clustername).
		WithField("task", fmt.Sprintf("pre-create cluster check"))

	clusterKeyPath := resource.NewCluster(clustername).BuildKey()

	return &PreCreateCheck{
		clustername:    clustername,
		clusterKeyPath: clusterKeyPath,
		log:            preCreateCheckLog,
		store:          store,
	}
}

func (p *PreCreateCheck) Run() error {
	clusterExists := p.store.KeyExists(p.clusterKeyPath)
	if clusterExists {
		return fmt.Errorf("cluster %s already exists", p.clustername)
	}

	return nil
}

type CheckClusterNetworkNodes struct {
	clusterNetwork *network.ClusterNetwork
	log            logrus.FieldLogger
	pingChecker    utils.PingChecker
}

func NewCheckClusterNetworkNodes(clusterNetwork *network.ClusterNetwork, log logrus.FieldLogger, pc utils.PingChecker) *CheckClusterNetworkNodes {
	checkClusterNetworkNodesLog := log.WithField("task", fmt.Sprintf("check cluster network nodes in %v", clusterNetwork))

	return &CheckClusterNetworkNodes{
		clusterNetwork: clusterNetwork,
		log:            checkClusterNetworkNodesLog,
		pingChecker:    pc,
	}
}

func (c *CheckClusterNetworkNodes) Run() error {
	// Blocking check waiting for cluster gateway to be up
	gwReady := c.pingChecker.ICMPPingNWithDelay(c.clusterNetwork.Gateway.String(), config.NodePingTimeout, c.log,
		config.NodePingRetry, config.NodePingDelay)

	// Cluster gateway still not ready after max retry
	if !gwReady {
		return fmt.Errorf("timeout waiting for cluster gateway to be up '%v'", c.clusterNetwork)
	}

	return nil
}

type CheckNode struct {
	ip          net.IP
	log         logrus.FieldLogger
	pingChecker utils.PingChecker
}

func NewCheckNode(ip net.IP, log logrus.FieldLogger, pc utils.PingChecker) *CheckNode {
	checkNodeLog := log.WithField("task", fmt.Sprintf("check node with address %v", ip.String()))

	return &CheckNode{
		ip:          ip,
		log:         checkNodeLog,
		pingChecker: pc,
	}
}

func (c *CheckNode) Run() error {
	// Blocking check waiting for node to be up
	nodeReady := c.pingChecker.ICMPPingNWithDelay(c.ip.String(), config.NodePingTimeout, c.log,
		config.NodePingRetry, config.NodePingDelay)

	// Cluster gateway still not ready after max retry
	if !nodeReady {
		return fmt.Errorf("timeout waiting for node to be up '%v'", c.ip)
	}

	return nil
}

type CreateNode struct {
	Config node.NodeCreateConfig
	p      Provisioner
	log    logrus.FieldLogger
}

func NewCreateNode(config node.NodeCreateConfig, p Provisioner, log logrus.FieldLogger) *CreateNode {
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
		n.log.Errorf("error creating node '%v': %s", n.Config, err)
		return err
	}
	return nil
}

type CreateAdminKubeconfig struct {
	clustername      string
	apiServerAddress string
	apiServerPort    string
	utils.SCPConfig
	log         logrus.FieldLogger
	pingChecker utils.PingChecker
}

func NewCreateAdminKubeconfig(clustername, apiServerAddress, apiServerPort string,
	config utils.SCPConfig, pc utils.PingChecker) *CreateAdminKubeconfig {

	createAdminKubeconfigLog := logrus.WithField("cluster", clustername).
		WithField("task", fmt.Sprintf("copying kubeconfig from %s to local system", config.RemoteIPAddress))

	return &CreateAdminKubeconfig{
		clustername:      clustername,
		apiServerAddress: apiServerAddress,
		apiServerPort:    apiServerPort,
		SCPConfig:        config,
		log:              createAdminKubeconfigLog,
		pingChecker:      pc,
	}
}

func (k *CreateAdminKubeconfig) Run() error {
	// Blocking check waiting control plane to be up
	apiServer := fmt.Sprintf("%s:%s", k.apiServerAddress, k.apiServerPort)
	ready := k.pingChecker.PortPingNWithDelay(apiServer, config.NodePingTimeout, k.log,
		config.NodePingRetry, config.NodePingDelay)
	if !ready {
		return fmt.Errorf("timeout waiting for control plane to be up '%v'", k)
	}
	// TODO(giri): wait until kubeadm finishes bootstraping,
	// hardcoded wait time 120 sec in cloud init script
	time.Sleep(3 * config.NodePingDelay)

	err := os.MkdirAll(config.MahakamMultiKubeconfigPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating mahakam directory to hold kubeconfig files")
	}

	// Copy over kubeconfig generated by kubeadm to mahakam server
	s := utils.NewSCPClient()
	_, err = s.CopyRemoteFile(k.SCPConfig)
	if err != nil {
		return fmt.Errorf("error creating admin kubeconfig file for cluster %s '%v': %s", k.clustername, k.SCPConfig, err)
	}

	k.log.Infof("admin kubeconfig has been copied over successfully '%v'", k)
	return nil
}

type ClusterNetworkReachability struct {
	ipAssigner        utils.IPAssigner
	mahakamServerIP   net.IP
	mahakamServerMask net.IPMask
	mahakamNetif      string
	log               logrus.FieldLogger
}

func NewClusterNetworkReachability(ipAssigner utils.IPAssigner, mahakamServerIP net.IP, mahakamServerMask net.IPMask, mahakamNetif string) *ClusterNetworkReachability {
	networkReachbilityLog := logrus.WithField("task", fmt.Sprintf("configure reachability from server to cluster network with %s and mask %s on %s", mahakamServerIP, mahakamServerMask, mahakamNetif))

	return &ClusterNetworkReachability{
		ipAssigner:        ipAssigner,
		mahakamServerIP:   mahakamServerIP,
		mahakamServerMask: mahakamServerMask,
		mahakamNetif:      mahakamNetif,
		log:               networkReachbilityLog,
	}
}

func (c *ClusterNetworkReachability) Run() error {
	_, err := c.ipAssigner.Assign(c.mahakamServerIP, c.mahakamServerMask, c.mahakamNetif)
	if err != nil {
		return fmt.Errorf("error assigning cluster network IP to mahakam server: %s", err)
	}
	return nil
}

type ClusterValidation struct {
	owner            string
	clustername      string
	clusterValidator validation.Validator
	store            store.ResourceStore
	log              logrus.FieldLogger
}

func NewClusterValidation(owner, clustername string, cv validation.Validator, s store.ResourceStore, log logrus.FieldLogger) *ClusterValidation {
	clusterValidationLog := log.WithField("cluster", clustername).
		WithField("task", fmt.Sprintf("validate cluster is ready and healthy"))

	return &ClusterValidation{
		owner:            owner,
		clustername:      clustername,
		clusterValidator: cv,
		store:            s,
		log:              clusterValidationLog,
	}
}

func (v *ClusterValidation) Run() error {
	// Blocking waiting cluster to be healthy
	ready := v.clusterValidator.ValidateNWithDelay(v.owner, v.clustername, config.NodePingTimeout, v.log,
		config.NodePingRetry, config.NodePingDelay)
	if !ready {
		return fmt.Errorf("timeout waiting for cluster %s to be ready", v.clustername)
	}

	if err := v.updateClusterResource(v.clustername); err != nil {
		return fmt.Errorf("error updating cluster resource with ready status: %s", err)
	}

	return nil
}

func (v *ClusterValidation) updateClusterResource(clustername string) error {
	c := resource.NewCluster(clustername)
	if err := v.store.Get(c); err != nil {
		return fmt.Errorf("error adding cluster resource into kv store '%v': %s", c, err)
	}

	c.Status = resource.StatusReady
	if _, err := v.store.Update(c); err != nil {
		return fmt.Errorf("error updating cluster resource into kv store '%v': %s", c, err)
	}

	return nil
}
