package agent

import (
	"time"

	"github.com/digitalocean/go-openvswitch/ovs"
	"github.com/sirupsen/logrus"

	"github.com/mahakamcloud/mahakam/pkg/api/v1"
	mahakamclient "github.com/mahakamcloud/mahakam/pkg/client"
	"github.com/mahakamcloud/mahakam/pkg/netd/util"
)

const (
	// defaultDelay configures delay between API server poll in seconds
	defaultDelay = 5
)

type provisionAgent struct {
	hostAddress   string
	netReconciler Reconciler

	ovsClient     *ovs.Client
	mahakamClient v1.ClusterAPI

	log logrus.FieldLogger
}

func NewProvisionAgent(hostAddress, mahakamAPIServer string, log logrus.FieldLogger) Agent {
	mahakamClient := mahakamclient.GetMahakamClusterClient(mahakamAPIServer)
	ovsClient := ovs.New()

	netReconciler := NewNetworkReconciler(mahakamClient, ovsClient)

	paLog := log.WithField("agent", "provision")

	return &provisionAgent{
		netReconciler: netReconciler,
		log:           paLog,
	}
}

func (pa *provisionAgent) Run() {
	pa.log.Info("Starting provision agent...")

	for {
		delay := defaultDelay + util.RandomIntRange(1, 5)
		time.Sleep(time.Duration(delay) * time.Second)

		if err := pa.Execute(); err != nil {
			pa.log.Errorf("error provisioning networks: %s", err)
		}
	}
}

func (pa *provisionAgent) Execute() error {
	var err error

	expectedClusters, err := pa.netReconciler.GetExpected()
	if err != nil {
		return err
	}

	actualClusterInfras, err := pa.netReconciler.GetActual()
	if err != nil {
		return err
	}

	// Reconcile between list of clusters from API server (desired state)
	// and list of bridges/tunnels from OVS daemon (actual state). Perform
	// action as necessary.
	states, err := pa.netReconciler.Reconcile(expectedClusters, actualClusterInfras)
	if err != nil {
		return err
	}

	for _, st := range states {
		switch st.action {
		case actionCreate:
			err = pa.provisionClusterHostGRE()
		}
	}

	return err
}

func (pa *provisionAgent) provisionClusterHostGRE() error {
	return nil
}
