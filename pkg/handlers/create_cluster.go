package handlers

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/clusters"
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/network"
	"github.com/mahakamcloud/mahakam/pkg/node"
	"github.com/mahakamcloud/mahakam/pkg/provisioner"
	r "github.com/mahakamcloud/mahakam/pkg/resource_store/resource"
	log "github.com/sirupsen/logrus"
)

// CreateCluster is handlers for create-cluster operation
type CreateCluster struct {
	Handlers
}

// Handle is handler for create-cluster operation
func (h *CreateCluster) Handle(params clusters.CreateClusterParams) middleware.Responder {

	// TODO(giri/vijay): create network first before creating cluster
	nwf, err := newCreateNetworkWF(params.Body, h.Handlers)
	if err != nil {
		log.Errorf("error creating network workflow %s", err)
		return clusters.NewCreateClusterDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: swag.String(fmt.Sprintf("error creating network workflow %s", err)),
		})
	}
	nwf.Run()

	c := r.NewResourceCluster(swag.StringValue(params.Body.Name))
	c.NumNodes = int(params.Body.NumNodes)
	c.Status = r.StatusPending
	c.NetworkName = nwf.clusterNetwork.Name

	// TODO(giri): create cluster workflow should pull
	// /etc/kubernetes/admin.conf into this kubeconfig path
	c.KubeconfigPath = h.generateKubeconfigPath(c.Owner, c.Name)

	_, err = h.Handlers.Store.Add(c)
	if err != nil {
		log.Errorf("error adding network resource into kv store '%v': %s", c, err)
		return clusters.NewCreateClusterDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: swag.String(err.Error()),
		})
	}

	wf, err := newCreateClusterWF(params.Body, nwf.clusterNetwork, h.Handlers)
	if err != nil {
		log.Errorf("error creating cluster workflow %s", err)
		return clusters.NewCreateClusterDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: swag.String(fmt.Sprintf("error creating cluster workflow %s", err)),
		})
	}
	wf.Run()

	// TODO(giri/iqbal): run this provisioner from another routine,
	// must update resource status to creating and success accordingly
	err = provisioner.CreateCluster(params.Body)
	if err != nil {
		log.Errorf("error provisioning cluster creation %s", err)
		return clusters.NewCreateClusterDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: swag.String(err.Error()),
		})
	}
	return clusters.NewCreateClusterCreated().WithPayload(params.Body)
}

func (h *CreateCluster) generateKubeconfigPath(owner, clusterName string) string {
	return fmt.Sprintf(config.MahakamMultiKubeconfigPath + "/" + owner + "-" + clusterName + "-kubeconfig")
}

// TODO(giri/vijay): create network
type createNetworkWF struct {
	handlers       Handlers
	clusterNetwork *network.ClusterNetwork
}

func newCreateNetworkWF(cluster *models.Cluster, handlers Handlers) (*createNetworkWF, error) {
	return &createNetworkWF{
		handlers: handlers,
	}, nil
}

func (cn *createNetworkWF) Run() error {
	// TODO(giri/vijay): create network
	n, err := cn.handlers.Network.AllocateClusterNetwork()
	cn.clusterNetwork = n
	log.Infof("cluster network has been created %v", cn.clusterNetwork)
	return err
}

type createClusterWF struct {
	handlers       Handlers
	clusterNetwork *network.ClusterNetwork
	controlPlane   node.Node
	workers        []node.Node

	controlPlaneNetworkConfig *node.NetworkConfig
}

func newCreateClusterWF(cluster *models.Cluster, clusterNetwork *network.ClusterNetwork, handlers Handlers) (*createClusterWF, error) {
	clusterName := swag.StringValue(cluster.Name)

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

	return &createClusterWF{
		handlers:       handlers,
		clusterNetwork: clusterNetwork,
		controlPlane:   controlPlane,
		workers:        workers,
	}, nil
}

func (c *createClusterWF) Run() error {
	tasks, err := c.getCreateTask()
	if err != nil {
		return err
	}

	go func() {
		for _, task := range tasks {
			if err := task.Run(); err != nil {
				log.Errorf("error running task %v: %s", task, err)
			}
		}
	}()
	return nil
}

func (c *createClusterWF) getCreateTask() ([]provisioner.Task, error) {
	var tasks []provisioner.Task
	tasks = c.setupControlPlaneSteps(tasks)
	tasks = c.setupWorkerSteps(tasks)
	return tasks, nil
}

func (c *createClusterWF) setupControlPlaneSteps(tasks []provisioner.Task) []provisioner.Task {
	return tasks
}

func (c *createClusterWF) setupWorkerSteps(tasks []provisioner.Task) []provisioner.Task {
	return tasks
}

func getNetworkConfig(clusterNetwork *network.ClusterNetwork) (*node.NetworkConfig, error) {
	macAddress := network.GenerateMacAddress()
	ip, err := clusterNetwork.AllocateIP()
	if err != nil {
		return nil, err
	}

	networkConfig := &node.NetworkConfig{
		MacAddress: macAddress,
		IP:         []byte(ip),
		Gateway:    clusterNetwork.Gateway,
		Nameserver: clusterNetwork.Nameserver,
	}
	return networkConfig, nil
}
