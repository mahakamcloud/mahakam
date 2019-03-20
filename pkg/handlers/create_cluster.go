package handlers

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/mahakamcloud/mahakam/pkg/scheduler"
	"github.com/mahakamcloud/mahakam/pkg/validation"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/client"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/client/networks"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/clusters"
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/network"
	"github.com/mahakamcloud/mahakam/pkg/provisioner"
	r "github.com/mahakamcloud/mahakam/pkg/resource_store/resource"
	"github.com/mahakamcloud/mahakam/pkg/task"
	"github.com/mahakamcloud/mahakam/pkg/utils"
	"github.com/sirupsen/logrus"
)

// CreateCluster is handlers for create-cluster operation
type CreateCluster struct {
	Handlers
	KubernetesConfig config.KubernetesConfig
	hosts            []config.Host
	log              logrus.FieldLogger
}

// NewCreateClusterHandler creates new CreateCluster object
func NewCreateClusterHandler(handlers Handlers) *CreateCluster {
	return &CreateCluster{
		Handlers:         handlers,
		KubernetesConfig: handlers.AppConfig.KubernetesConfig,
		hosts:            handlers.AppConfig.HostsConfig,
		log:              handlers.Log,
	}
}

// Handle is handler for create-cluster operation
func (h *CreateCluster) Handle(params clusters.CreateClusterParams) middleware.Responder {
	h.log.Infof("handling create cluster request: %v", params)

	pwf, err := newPreCreateClusterWF(params.Body, h)
	if err != nil {
		h.log.Errorf("error creating pre-create cluster workflow %s", err)
		return clusters.NewCreateClusterDefault(http.StatusInternalServerError).WithPayload(&models.Error{
			Message: fmt.Sprintf("error creating pre-check workflow %s", err),
		})
	}
	// blocking pre-check run
	err = pwf.Run()
	if err != nil {
		return clusters.NewCreateClusterDefault(http.StatusBadRequest).WithPayload(&models.Error{
			Message: fmt.Sprintf("error running pre-create cluster workflow %s", err),
		})
	}
	h.log.Debugf("pre-create cluster workflow completes successfully: %v", params.Body)

	// TODO(giri): must update resource status to creating and succcess accordingly
	wf, err := newCreateClusterWF(params.Body, h)
	if err != nil {
		h.log.Errorf("error creating cluster workflow %s", err)
		return clusters.NewCreateClusterDefault(http.StatusBadRequest).WithPayload(&models.Error{
			Message: fmt.Sprintf("error creating cluster workflow %s", err),
		})
	}
	wf.Run()

	cres := &models.Cluster{
		Name:     params.Body.Name,
		NumNodes: params.Body.NumNodes,
		NodeSize: params.Body.NodeSize,
		Status:   string(r.StatusPending),
	}
	return clusters.NewCreateClusterCreated().WithPayload(cres)
}

type preCreateClusterWF struct {
	handlers    Handlers
	log         logrus.FieldLogger
	owner       string
	clustername string
}

func newPreCreateClusterWF(cluster *models.Cluster, cHandler *CreateCluster) (*preCreateClusterWF, error) {
	pwfLog := cHandler.log.WithField("workflow", "pre-check")
	pwfLog.Debugf("init pre-check workflow")

	clusterName := swag.StringValue(cluster.Name)

	return &preCreateClusterWF{
		handlers:    cHandler.Handlers,
		log:         pwfLog,
		owner:       cluster.Owner,
		clustername: clusterName,
	}, nil
}

func (p *preCreateClusterWF) Run() error {
	p.log.Infof("running pre-create cluster workflow: %v", p)

	tasks, err := p.getPreCreateTask()
	if err != nil {
		return err
	}

	for _, t := range tasks {
		p.log.Infof("running task %v", t)
		if err := t.Run(); err != nil {
			p.log.Errorf("error running task %v: %s", t, err)
			return fmt.Errorf("error pre-create cluster workflow %s", err)
		}
	}
	return nil
}

func (p *preCreateClusterWF) getPreCreateTask() ([]task.Task, error) {
	p.log.Debugf("getting pre-create task for cluster %s", p.clustername)

	var tasks []task.Task
	tasks = p.setupPreCheckTasks(tasks)
	return tasks, nil
}

func (p *preCreateClusterWF) setupPreCheckTasks(tasks []task.Task) []task.Task {
	p.log.Debugf("setup pre-check tasks for cluster %s", p.clustername)

	preCreateCheck := provisioner.NewPreCreateCheck(p.clustername, p.log, p.handlers.Store)

	tasks = append(tasks, preCreateCheck)
	return tasks
}

type createClusterWF struct {
	handlers       Handlers
	log            logrus.FieldLogger
	owner          string
	clustername    string
	clusterNetwork *network.ClusterNetwork
	controlPlane   r.Node
	workers        []r.Node
	workerCount    int
	hosts          []config.Host

	controlPlaneIP net.IP
	nodePublicKey  string
	podNetworkCidr string
	kubeadmToken   string
}

func newCreateClusterWF(cluster *models.Cluster, cHandler *CreateCluster) (*createClusterWF, error) {
	cwfLog := cHandler.log.WithField("workflow", "create-cluster")
	cwfLog.Debugf("init create cluster workflow")

	clusterName := swag.StringValue(cluster.Name)
	workerNodeSize := swag.StringValue(cluster.NodeSize)

	if workerNodeSize == "" {
		workerNodeSize = r.ClusterSizeDefault
	}

	if !r.ClusterSizeValidate(workerNodeSize) {
		cwfLog.Errorf("provided cluster size is not available: %s", workerNodeSize)
		return nil, fmt.Errorf("provided cluster size is not available: %s", workerNodeSize)
	}

	workerNumCPUs := r.GetClusterNodeCPUs(workerNodeSize)
	workerMemorySize, err := r.GetClusterNodeMemoryInMB(workerNodeSize)
	if err != nil {
		cwfLog.Errorf("error getting memory size %s", err)
		return nil, err
	}

	// For controlplane default to ClusterSizeDefault
	// instead of NodeSize passed from CLI
	cpNodeSize := r.ClusterSizeDefault
	cpNumCPUs := r.GetClusterNodeCPUs(cpNodeSize)

	cpMemorySize, err := r.GetClusterNodeMemoryInMB(cpNodeSize)
	if err != nil {
		cwfLog.Errorf("error getting memory size %s", err)
		return nil, err
	}

	clusterNetwork, err := getClusterNetwork(clusterName, cHandler.Network, cHandler.Client)
	if err != nil {
		cwfLog.Errorf("error getting network allocation for cluster %s: %s", clusterName, err)
	}

	workerNodesCount := int(cluster.NumNodes)

	err = storeClusterResource(clusterName, workerNodesCount, workerNodeSize, clusterNetwork, cHandler)
	if err != nil {
		cwfLog.Errorf("error storing cluster resource to kvstore '%v': %s", cluster, err)
		return nil, err
	}

	return &createClusterWF{
		handlers:       cHandler.Handlers,
		log:            cwfLog,
		owner:          cluster.Owner,
		clustername:    clusterName,
		clusterNetwork: clusterNetwork,
		controlPlane:   r.Node{},
		workers:        []r.Node{},
		hosts:          cHandler.hosts,
		controlPlaneIP: nil,
		workerCount:    workerNodesCount,
		nodePublicKey:  cHandler.KubernetesConfig.SSHPublicKey,
		podNetworkCidr: cHandler.KubernetesConfig.PodNetworkCidr,
		kubeadmToken:   cHandler.KubernetesConfig.KubeadmToken,
	}, nil
}

func (c *createClusterWF) Run() error {
	c.log.Infof("running create cluster workflow: %v", c)

	tasks, err := c.getCreateTask()
	if err != nil {
		return err
	}

	for i := range tasks {
		go func(t task.Task) {
			c.log.Infof("running task %v", t)
			if err := t.Run(); err != nil {
				c.log.Errorf("error running task %v: %s", t, err)
			}
		}(tasks[i])
	}
	return nil
}

func (c *createClusterWF) getCreateTask() ([]task.Task, error) {
	c.log.Debugf("getting create task for cluster %s", c.clustername)

	var tasks []task.Task
	tasks = c.setupControlPlaneTasks(tasks)
	tasks = c.setupWorkerTasks(tasks)
	tasks = c.setupAdminKubeconfigTasks(tasks)
	tasks = c.setupClusterValidationTasks(tasks)
	return tasks, nil
}

func (c *createClusterWF) setupControlPlaneTasks(tasks []task.Task) []task.Task {
	c.log.Debugf("setup control plane tasks for cluster %s", c.clustername)

	host, err := scheduler.GetHost(c.hosts)
	if err != nil {
		c.log.Errorf("Error : %v", err)
		return nil
	}

	cpNodeSpec := r.NewNodeSpec(r.SmallNode)
	cpNodeStatus := r.NewNodeStatus(host.String())
	cpNodeFqdn := getFQDN(c.controlPlane.Name, c.clusterNetwork.Name, c.handlers.AppConfig.NetworkConfig.Domain)

	extraConfig := map[string]string{
		config.KeyPodNetworkCidr: c.podNetworkCidr,
		config.KeyKubeadmToken:   c.kubeadmToken,
	}
	cpNodeMetadata := r.NewMetadata("", []string{c.nodePublicKey}, extraConfig)

	cpNetworkConfigs := []r.NetworkConfig{
		*r.NewNetworkConfig(network.GenerateMacAddress(),
			cpNodeFqdn,
			c.clusterNetwork.ClusterNetworkCIDR.Mask,
			nil,
			c.clusterNetwork.Gateway,
			nil),
	}

	cpNodeConfig := r.NewNode(fmt.Sprintf("%s-cp", c.clustername),
		*cpNodeSpec,
		cpNetworkConfigs,
		*cpNodeMetadata,
		*cpNodeStatus,
		string(r.RoleNetworkGW))

	checkClusterNetworkNodes := provisioner.NewCheckClusterNetworkNodes(c.clusterNetwork, c.log, utils.NewPingCheck())
	createControlPlaneNode := provisioner.NewCreateNode(*cpNodeConfig, c.handlers.Provisioner, c.log)

	controlPlaneSeqTasks := task.NewSeqTask(c.log, checkClusterNetworkNodes, createControlPlaneNode)
	tasks = append(tasks, controlPlaneSeqTasks)

	return tasks
}

func (c *createClusterWF) setupWorkerTasks(tasks []task.Task) []task.Task {
	c.log.Debugf("setup worker steps for cluster %s", c.clustername)

	for i := 1; i <= c.workerCount; i++ {
		host, err := scheduler.GetHost(c.hosts)
		if err != nil {
			c.log.Errorf("Error : %v", err)
			return nil
		}

		wrNodeSpec := r.NewNodeSpec(r.SmallNode)
		wrNodeStatus := r.NewNodeStatus(host.String())
		wrNodeFqdn := getFQDN(c.controlPlane.Name, c.clusterNetwork.Name, c.handlers.AppConfig.NetworkConfig.Domain)

		// TODO: We don't know control plane IP here, due to DHCP allocation.
		// This needs to be sorted either via pre allocated reserved IP or
		// some other mechanism. This breaks cluster creation workflow
		extraConfig := map[string]string{
			config.KeyControlPlaneIP: c.controlPlaneIP.String(),
			config.KeyKubeadmToken:   c.kubeadmToken,
		}
		wrNodeMetadata := r.NewMetadata("", []string{c.nodePublicKey}, extraConfig)

		wrNetworkConfigs := []r.NetworkConfig{
			*r.NewNetworkConfig(network.GenerateMacAddress(),
				wrNodeFqdn,
				c.clusterNetwork.ClusterNetworkCIDR.Mask,
				nil,
				c.clusterNetwork.Gateway,
				nil),
		}

		wrNodeConfig := r.NewNode(fmt.Sprintf("%s-worker-%d", c.clustername, i),
			*wrNodeSpec,
			wrNetworkConfigs,
			*wrNodeMetadata,
			*wrNodeStatus,
			string(r.RoleNetworkGW))

		checkClusterNetworkNodes := provisioner.NewCheckClusterNetworkNodes(c.clusterNetwork, c.log, utils.NewPingCheck())
		createWorkerNode := provisioner.NewCreateNode(*wrNodeConfig, c.handlers.Provisioner, c.log)

		workerSeqTasks := task.NewSeqTask(c.log, checkClusterNetworkNodes, createWorkerNode)

		tasks = append(tasks, workerSeqTasks)
	}

	return tasks
}

func (c *createClusterWF) setupAdminKubeconfigTasks(tasks []task.Task) []task.Task {
	c.log.Debugf("setup admin kubeconfig steps for cluster %s", c.clustername)

	sConfig := utils.SCPConfig{
		Username:        config.KubernetesNodeUsername,
		RemoteIPAddress: c.controlPlaneIP.String(),
		PrivateKeyPath:  config.MahakamSSHPrivateKeyPath,
		RemoteFilePath:  config.KubernetesAdminKubeconfigPath,
		LocalFilePath: utils.GenerateKubeconfigPath(
			config.MahakamMultiKubeconfigPath,
			c.owner,
			c.clustername,
		),
	}

	createAdminKubeconfig := provisioner.NewCreateAdminKubeconfig(c.clustername,
		c.controlPlaneIP.String(), strconv.Itoa(config.KubernetesAPIServerPort), sConfig, utils.NewPingCheck())

	tasks = append(tasks, createAdminKubeconfig)
	return tasks
}

func (c *createClusterWF) setupClusterValidationTasks(tasks []task.Task) []task.Task {
	c.log.Debugf("setup cluster validator steps for cluster %s", c.clustername)

	clusterValidation := provisioner.NewClusterValidation(c.owner, c.clustername, validation.NewClusterValidator(c.handlers.Client), c.handlers.Store, c.log)

	tasks = append(tasks, clusterValidation)
	return tasks
}

func storeClusterResource(clusterName string, numNodes int, nodeSize string, clusternet *network.ClusterNetwork, cHandler *CreateCluster) error {
	c := r.NewCluster(clusterName)
	c.NumNodes = numNodes
	c.NodeSize = nodeSize
	c.Status = r.StatusPending
	c.NetworkName = clusternet.Name
	c.KubeconfigPath = utils.GenerateKubeconfigPath(config.MahakamMultiKubeconfigPath, c.Owner, c.Name)

	_, err := cHandler.Handlers.Store.Add(c)
	if err != nil {
		return fmt.Errorf("error adding cluster resource into kv store '%v': %s", c, err)
	}
	return nil
}

func getClusterNetwork(clusterName string, netmanager *network.NetworkManager, c *client.Mahakam) (*network.ClusterNetwork, error) {
	req := &models.Network{
		Name: swag.String(clusterName),
	}

	res, err := c.Networks.CreateNetwork(networks.NewCreateNetworkParams().WithBody(req))
	if err != nil {
		return nil, err
	}

	_, cidr, _ := net.ParseCIDR(res.Payload.NetworkCIDR)
	return network.NewClusterNetwork(*cidr, netmanager), nil
}

func getFQDN(hostname, clusterNetworkName, rootDomain string) string {
	return hostname + "." + clusterNetworkName + "." + rootDomain
}
