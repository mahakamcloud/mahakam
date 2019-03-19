package agent

import (
	"github.com/digitalocean/go-openvswitch/ovs"
	"github.com/mahakamcloud/mahakam/pkg/api/v1"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/client/clusters"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/netd/network"
)

type reconcileAction string

const (
	actionCreate  reconcileAction = "create"
	actionDestroy reconcileAction = "destroy"
	actionNone    reconcileAction = "none"
)

type ReconcileState struct {
	clusterName string
	action      reconcileAction
}

type Reconciler interface {
	GetExpected() ([]*models.Cluster, error)
	GetActual() ([]*network.ClusterHostGRE, error)
	Reconcile(expected []*models.Cluster, actual []*network.ClusterHostGRE) ([]*ReconcileState, error)
}

type NetworkReconciler struct {
	mahakamClient v1.ClusterAPI
	ovsClient     *ovs.Client
}

func NewNetworkReconciler(mahakamClient v1.ClusterAPI, ovsClient *ovs.Client) Reconciler {
	return &NetworkReconciler{
		mahakamClient: mahakamClient,
		ovsClient:     ovsClient,
	}
}

func (nr *NetworkReconciler) GetExpected() ([]*models.Cluster, error) {
	res, err := nr.mahakamClient.GetClusters(clusters.NewGetClustersParams())
	if err != nil {
		return nil, err
	}
	return res.Payload, nil
}

func (nr *NetworkReconciler) GetActual() ([]*network.ClusterHostGRE, error) {
	// TODO(giri): provision GRE tunnel with OVS client
	return nil, nil
}

func (nr *NetworkReconciler) Reconcile(
	expectedClusters []*models.Cluster,
	actualClusterHostGREs []*network.ClusterHostGRE) ([]*ReconcileState, error) {
	// TODO(giri): reconcile states
	return nil, nil
}
