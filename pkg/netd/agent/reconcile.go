package agent

import (
	"github.com/digitalocean/go-openvswitch/ovs"
	"github.com/go-openapi/swag"

	"github.com/mahakamcloud/mahakam/pkg/api/v1"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/client/clusters"
	"github.com/mahakamcloud/mahakam/pkg/netd/network"
)

type reconcileAction string

const (
	actionCreate  reconcileAction = "create"
	actionDestroy reconcileAction = "destroy"
	actionNone    reconcileAction = "none"
)

type State struct {
	network.ClusterHostGRE
	action reconcileAction
}

type ReconcileStates struct {
	states map[string]*State
}

type Reconciler interface {
	GetExpected() (*ReconcileStates, error)
	GetActual() (*ReconcileStates, error)
	Reconcile(expected *ReconcileStates, actual *ReconcileStates) (*ReconcileStates, error)
}

type NetworkReconciler struct {
	mahakamClient v1.ClusterAPI
	ovsClient     *ovs.Client

	reconciledStates *ReconcileStates
}

func NewNetworkReconciler(mahakamClient v1.ClusterAPI, ovsClient *ovs.Client) Reconciler {
	return &NetworkReconciler{
		mahakamClient: mahakamClient,
		ovsClient:     ovsClient,
	}
}

func (nr *NetworkReconciler) GetExpected() (*ReconcileStates, error) {
	// TODO(giri): get GRE key from API server
	res, err := nr.mahakamClient.GetClusters(clusters.NewGetClustersParams())
	if err != nil {
		return nil, err
	}

	expectedClusterStates := make(map[string]*State)
	clusters := res.Payload
	for _, cl := range clusters {
		// TODO(giri): populate proper GRE key
		key := swag.StringValue(cl.Name)
		expectedClusterStates[key] = &State{
			ClusterHostGRE: network.ClusterHostGRE{
				GREKey: "1",
			},
		}
	}
	return &ReconcileStates{
		states: expectedClusterStates,
	}, nil
}

func (nr *NetworkReconciler) GetActual() (*ReconcileStates, error) {
	// TODO(giri): gather GRE tunnel info with OVS client
	actualClusterStates := make(map[string]*State)

	return &ReconcileStates{
		states: actualClusterStates,
	}, nil
}

func (nr *NetworkReconciler) Reconcile(
	expectedClusterStates *ReconcileStates,
	actualClusterStates *ReconcileStates) (*ReconcileStates, error) {

	reconciledStates := make(map[string]*State)

	for cl := range expectedClusterStates.states {
		if _, ok := actualClusterStates.states[cl]; !ok {
			reconciledStates[cl] = &State{
				ClusterHostGRE: expectedClusterStates.states[cl].ClusterHostGRE,
				action:         actionCreate,
			}
		}
	}

	for cl := range actualClusterStates.states {
		if _, ok := expectedClusterStates.states[cl]; !ok {
			reconciledStates[cl] = &State{
				ClusterHostGRE: actualClusterStates.states[cl].ClusterHostGRE,
				action:         actionDestroy,
			}
		}
	}

	return &ReconcileStates{
		states: reconciledStates,
	}, nil
}
