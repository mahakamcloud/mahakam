package handlers

import (
	"fmt"
	"net"
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/client/networks"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/clusters"
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/network"
	"github.com/mahakamcloud/mahakam/pkg/node"
	"github.com/mahakamcloud/mahakam/pkg/provisioner"
	r "github.com/mahakamcloud/mahakam/pkg/resource_store/resource"
	"github.com/mahakamcloud/mahakam/pkg/task"
	"github.com/mahakamcloud/mahakam/pkg/utils"
	log "github.com/sirupsen/logrus"
)

// CreateCluster is handlers for create-cluster operation
type CreateCluster struct {
	Handlers
	KubernetesConfig config.KubernetesConfig
	log              log.FieldLogger
}

func NewCreateClusterHandler(handlers Handlers, config config.KubernetesConfig) *CreateCluster {
	log := log.WithField("create", "cluster")
	return &CreateCluster{
		Handlers:         handlers,
		KubernetesConfig: config,
		log:              log,
	}
}

// Handle is handler for create-cluster operation
func (h *CreateCluster) Handle(params clusters.CreateClusterParams) middleware.Responder {

	// TODO(giri): must update resource status to creating and succcess accordingly
	wf, err := newCreateClusterWF(params.Body, h)
	if err != nil {
		log.Errorf("error creating cluster workflow %s", err)
		return clusters.NewCreateClusterDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: fmt.Sprintf("error creating cluster workflow %s", err),
		})
	}
	wf.Run()

	cres := &models.Cluster{
		Name:        params.Body.Name,
		ClusterPlan: params.Body.ClusterPlan,
		NumNodes:    params.Body.NumNodes,
		Status:      string(r.StatusPending),
	}
	return clusters.NewCreateClusterCreated().WithPayload(cres)
}

type createClusterWF struct {
	handlers       Handlers
	log            log.FieldLogger
	owner          string
	clustername    string
	clusterNetwork *network.ClusterNetwork
	controlPlane   node.Node
	workers        []node.Node

	controlPlaneIP net.IP
	nodePublicKey  string
	podNetworkCidr string
	kubeadmToken   string
}

func newCreateClusterWF(cluster *models.Cluster, cHandler *CreateCluster) (*createClusterWF, error) {
	clusterName := swag.StringValue(cluster.Name)

	clusterNetwork, err := getClusterNetwork(clusterName, cHandler.Network)
	if err != nil {
		log.Errorf("error getting network allocation for cluster %s: %s", clusterName, err)
	}

	controlPlaneNetworkConfig, err := getNetworkConfig(clusterNetwork)
	if err != nil {
		log.Errorf("error getting network config for control plane: %s", err)
		return nil, err
	}

	controlPlane := node.Node{
		Name:          fmt.Sprintf("%s-cp", clusterName),
		NetworkConfig: *controlPlaneNetworkConfig,
	}

	var workers []node.Node
	workerNodesCount := int(cluster.NumNodes)
	for i := 1; i <= workerNodesCount; i++ {
		workerNetworkConfig, err := getNetworkConfig(clusterNetwork)
		if err != nil {
			return nil, err
		}

		worker := node.Node{
			Name:          fmt.Sprintf("%s-worker-%d", clusterName, i),
			NetworkConfig: *workerNetworkConfig,
		}
		workers = append(workers, worker)
	}

	err = storeClusterResource(clusterName, workerNodesCount, clusterNetwork, cHandler)
	if err != nil {
		log.Errorf("error storing cluster resource to kvstore '%v': %s", cluster, err)
		return nil, err
	}

	return &createClusterWF{
		handlers:       cHandler.Handlers,
		log:            cHandler.log,
		owner:          cluster.Owner,
		clustername:    clusterName,
		clusterNetwork: clusterNetwork,
		controlPlane:   controlPlane,
		workers:        workers,
		controlPlaneIP: controlPlane.NetworkConfig.IP,
		nodePublicKey:  cHandler.KubernetesConfig.SSHPublicKey,
		podNetworkCidr: cHandler.KubernetesConfig.PodNetworkCidr,
		kubeadmToken:   cHandler.KubernetesConfig.KubeadmToken,
	}, nil
}

func (c *createClusterWF) Run() error {
	tasks, err := c.getCreateTask()
	if err != nil {
		return err
	}

	for i := range tasks {
		go func(t task.Task) {
			c.log.Infof("Running task %v", t)
			if err := t.Run(); err != nil {
				c.log.Errorf("error running task %v: %s", t, err)
			}
		}(tasks[i])
	}
	return nil
}

func (c *createClusterWF) getCreateTask() ([]task.Task, error) {
	var tasks []task.Task
	tasks = c.setupControlPlaneSteps(tasks)
	tasks = c.setupWorkerSteps(tasks)
	tasks = c.setupAdminKubeconfigSteps(tasks)
	return tasks, nil
}

func (c *createClusterWF) setupControlPlaneSteps(tasks []task.Task) []task.Task {
	cpConfig := node.NodeCreateConfig{
		// TODO(giri): must be getting from list of hosts
		Host: net.ParseIP("10.30.0.1"),
		Role: node.RoleControlPlane,
		Node: node.Node{
			Name:         c.controlPlane.Name,
			SSHPublicKey: c.nodePublicKey,
			NetworkConfig: node.NetworkConfig{
				MacAddress: c.controlPlane.MacAddress,
				IP:         c.controlPlane.IP,
				Mask:       c.controlPlane.Mask,
				Gateway:    c.controlPlane.Gateway,
				Nameserver: c.controlPlane.Nameserver,
			},
		},
		ExtraConfig: map[string]string{
			config.KeyPodNetworkCidr: c.podNetworkCidr,
			config.KeyKubeadmToken:   c.kubeadmToken,
		},
	}

	checkClusterNetworkNodes := provisioner.NewCheckClusterNetworkNodes(c.clusterNetwork, c.log)
	createControlPlaneNode := provisioner.NewCreateNode(cpConfig, c.handlers.Provisioner, c.log)

	controlPlaneSeqTasks := task.NewSeqTask(c.log, checkClusterNetworkNodes, createControlPlaneNode)
	tasks = append(tasks, controlPlaneSeqTasks)

	return tasks
}

func (c *createClusterWF) setupWorkerSteps(tasks []task.Task) []task.Task {
	for _, worker := range c.workers {
		wConfig := node.NodeCreateConfig{
			Host: net.ParseIP("10.30.0.1"),
			Role: node.RoleWorker,
			Node: node.Node{
				Name:         worker.Name,
				SSHPublicKey: c.nodePublicKey,
				NetworkConfig: node.NetworkConfig{
					MacAddress: worker.MacAddress,
					IP:         worker.IP,
					Mask:       worker.Mask,
					Gateway:    worker.Gateway,
					Nameserver: worker.Nameserver,
				},
			},
			ExtraConfig: map[string]string{
				config.KeyControlPlaneIP: c.controlPlaneIP.String(),
				config.KeyKubeadmToken:   c.kubeadmToken,
			},
		}

		checkClusterNetworkNodes := provisioner.NewCheckClusterNetworkNodes(c.clusterNetwork, c.log)
		createWorkerNode := provisioner.NewCreateNode(wConfig, c.handlers.Provisioner, c.log)

		workerSeqTasks := task.NewSeqTask(c.log, checkClusterNetworkNodes, createWorkerNode)

		tasks = append(tasks, workerSeqTasks)
	}

	return tasks
}

func (c *createClusterWF) setupAdminKubeconfigSteps(tasks []task.Task) []task.Task {
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
		c.controlPlaneIP.String(), strconv.Itoa(config.KubernetesAPIServerPort), sConfig)

	tasks = append(tasks, createAdminKubeconfig)
	return tasks
}

func storeClusterResource(clusterName string, numNodes int, clusternet *network.ClusterNetwork, cHandler *CreateCluster) error {
	c := r.NewResourceCluster(clusterName)
	c.NumNodes = numNodes
	c.Status = r.StatusPending
	c.NetworkName = clusternet.Name
	c.KubeconfigPath = utils.GenerateKubeconfigPath(config.MahakamMultiKubeconfigPath, c.Owner, c.Name)

	_, err := cHandler.Handlers.Store.Add(c)
	if err != nil {
		log.Errorf("error adding cluster resource into kv store '%v': %s", c, err)
		return fmt.Errorf("error adding cluster resource into kv store '%v': %s", c, err)
	}
	return nil
}

func getClusterNetwork(clusterName string, netmanager *network.NetworkManager) (*network.ClusterNetwork, error) {
	// TODO(giri): get local host and local port from config.yaml
	client := GetMahakamClient(":" + strconv.Itoa(config.MahakamAPIDefaultPort))
	req := &models.Network{
		Name: swag.String(clusterName),
	}

	res, err := client.Networks.CreateNetwork(networks.NewCreateNetworkParams().WithBody(req))
	if err != nil {
		return nil, err
	}

	_, cidr, _ := net.ParseCIDR(res.Payload.NetworkCIDR)
	return network.NewClusterNetwork(*cidr, netmanager), nil
}

func getNetworkConfig(clusterNetwork *network.ClusterNetwork) (*node.NetworkConfig, error) {
	macAddress := network.GenerateMacAddress()
	ip, err := clusterNetwork.AllocateIP()
	if err != nil {
		return nil, err
	}

	networkConfig := &node.NetworkConfig{
		MacAddress: macAddress,
		IP:         net.ParseIP(ip),
		Mask:       clusterNetwork.ClusterNetworkCIDR.Mask,
		Gateway:    clusterNetwork.Gateway,
		Nameserver: clusterNetwork.Nameserver,
	}
	return networkConfig, nil
}
