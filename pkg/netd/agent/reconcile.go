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
		greKey := "1"
		expectedClusterStates[greKey] = &State{
			ClusterHostGRE: network.ClusterHostGRE{
				ClusterName: swag.StringValue(cl.Name),
				GREKey:      greKey,
			},
		}
	}
	return &ReconcileStates{
		states: expectedClusterStates,
	}, nil
}

func (nr *NetworkReconciler) GetActual() (*ReconcileStates, error) {
	actualClusterStates := make(map[string]*State)

	brs, err := nr.ovsClient.VSwitch.ListBridges()
	if err != nil {
		return nil, err
	}

	// TODO(giri): currently this is only fetching bridge, remaining:
	// tap dev, number of tunnels
	for _, br := range brs {
		greKey := network.ParseGREKey(br)
		actualClusterStates[greKey] = &State{
			ClusterHostGRE: network.ClusterHostGRE{
				GREKey: greKey,
			},
		}
	}

	return &ReconcileStates{
		states: actualClusterStates,
	}, nil
}

func (nr *NetworkReconciler) Reconcile(
	expectedClusterStates *ReconcileStates,
	actualClusterStates *ReconcileStates) (*ReconcileStates, error) {

	reconciledStates := make(map[string]*State)

	for greKey := range expectedClusterStates.states {
		if _, ok := actualClusterStates.states[greKey]; !ok {
			reconciledStates[greKey] = &State{
				ClusterHostGRE: expectedClusterStates.states[greKey].ClusterHostGRE,
				action:         actionCreate,
			}
		}
	}

	for greKey := range actualClusterStates.states {
		if _, ok := expectedClusterStates.states[greKey]; !ok {
			reconciledStates[greKey] = &State{
				ClusterHostGRE: actualClusterStates.states[greKey].ClusterHostGRE,
				action:         actionDestroy,
			}
		}
	}

	return &ReconcileStates{
		states: reconciledStates,
	}, nil
}
