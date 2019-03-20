package agent

import (
	"fmt"
	"time"

	"github.com/digitalocean/go-openvswitch/ovs"
	"github.com/sirupsen/logrus"

	"github.com/mahakamcloud/mahakam/pkg/api/v1"
	mahakamclient "github.com/mahakamcloud/mahakam/pkg/client"
	"github.com/mahakamcloud/mahakam/pkg/netd/util"
	"github.com/mahakamcloud/mahakam/pkg/task"
)

const (
	// defaultDelay configures delay between API server poll in seconds
	defaultDelay = 5
)

type provisionAgent struct {
	hostname      string
	hostAddress   string
	netReconciler Reconciler

	ovsClient     *ovs.Client
	mahakamClient v1.ClusterAPI

	log logrus.FieldLogger
}

func NewProvisionAgent(clustername, hostname, hostAddress, mahakamAPIServer string, log logrus.FieldLogger) Agent {
	mahakamClient := mahakamclient.GetMahakamClusterClient(mahakamAPIServer)
	ovsClient := ovs.New()

	netReconciler := NewNetworkReconciler(mahakamClient, ovsClient)

	paLog := log.WithField("agent", "provision")

	return &provisionAgent{
		hostname:      hostname,
		hostAddress:   hostAddress,
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

	var provisionErrors []error
	for _, st := range states {
		switch st.action {
		case actionCreate:
			if err := pa.provisionClusterHostGRE(st.clustername); err != nil {
				provisionErrors = append(provisionErrors, fmt.Errorf("error provisioning %q cluster host gre: %s", st.clustername, err))
			}
		}
	}

	if len(provisionErrors) > 0 {
		for _, err := range provisionErrors {
			pa.log.Errorf("%s", err)
		}
		return fmt.Errorf("reconcile action errors out")
	}

	return nil
}

func (pa *provisionAgent) provisionClusterHostGRE(clustername string) error {
	return nil
}

func (pa *provisionAgent) getProvisionTask(clustername string) ([]task.Task, error) {
	pa.log.Debugf("getting provision task for cluster %q host gre", clustername)

	return nil, nil
}
