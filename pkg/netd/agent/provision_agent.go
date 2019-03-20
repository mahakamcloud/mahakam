package agent

import (
	"fmt"
	"net"
	"time"

	"github.com/digitalocean/go-openvswitch/ovs"
	"github.com/sirupsen/logrus"

	"github.com/mahakamcloud/mahakam/pkg/api/v1"
	mahakamclient "github.com/mahakamcloud/mahakam/pkg/client"
	"github.com/mahakamcloud/mahakam/pkg/netd/network"
	"github.com/mahakamcloud/mahakam/pkg/netd/provisioner"
	"github.com/mahakamcloud/mahakam/pkg/netd/util"
	"github.com/mahakamcloud/mahakam/pkg/task"
)

const (
	// defaultDelay configures delay between API server poll in seconds
	defaultDelay = 5
)

type provisionAgent struct {
	hostname    string
	localHostIP net.IP

	netReconciler Reconciler
	mahakamClient v1.ClusterAPI
	ovsClient     *ovs.Client

	log logrus.FieldLogger
}

func NewProvisionAgent(clustername, hostname, localIP, mahakamAPIServer string, log logrus.FieldLogger) Agent {
	mahakamClient := mahakamclient.GetMahakamClusterClient(mahakamAPIServer)
	ovsClient := ovs.New(ovs.Sudo())

	netReconciler := NewNetworkReconciler(mahakamClient, ovsClient)

	localHostIP := net.ParseIP(localIP)

	paLog := log.WithField("agent", "provision")

	return &provisionAgent{
		hostname:      hostname,
		localHostIP:   localHostIP,
		netReconciler: netReconciler,
		mahakamClient: mahakamClient,
		ovsClient:     ovsClient,
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
	// TODO(giri): getting list of hosts managed by Mahakam
	hostIPs := []net.IP{}

	// TODO(giri): cluster model must include GRE key
	expectedClusters, err := pa.netReconciler.GetExpected()
	if err != nil {
		return err
	}

	actualClusters, err := pa.netReconciler.GetActual()
	if err != nil {
		return err
	}

	// Reconcile between list of clusters from API server (expected state)
	// and list of bridges/tunnels from OVS daemon (actual state). Perform
	// action as necessary.
	reconciledClusters, err := pa.netReconciler.Reconcile(expectedClusters, actualClusters)
	if err != nil {
		return err
	}

	var provisionErrors []error
	for cl, val := range reconciledClusters.states {
		switch val.action {
		case actionCreate:
			grewf := newProvisionClusterHostGreWF(cl, val.GREKey, hostIPs, pa)
			if err := grewf.Run(); err != nil {
				provisionErrors = append(provisionErrors, fmt.Errorf("error provisioning %q cluster host gre: %s", cl, err))
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

type provisionClusterHostGreWF struct {
	clustername string
	*network.ClusterHostGRE
	*provisionAgent

	log logrus.FieldLogger
}

func newProvisionClusterHostGreWF(clustername, greKey string, hostIPs []net.IP, pa *provisionAgent) *provisionClusterHostGreWF {
	wfLog := pa.log.WithField("workflow", "provision-cluster-host-gre-wf")

	brName := fmt.Sprintf(network.BridgeFormat, greKey)

	var hostGRETunnels []network.HostGRETunnel
	for _, hostIP := range hostIPs {
		if hostIP.String() == pa.localHostIP.String() {
			continue
		}
		tunnel := network.HostGRETunnel{
			TapDevName:   fmt.Sprintf(network.TapDevFormat, greKey, util.LastOctet(pa.localHostIP), util.LastOctet(hostIP)),
			LocalHostIP:  pa.localHostIP,
			RemoteHostIP: hostIP,
		}
		hostGRETunnels = append(hostGRETunnels, tunnel)
	}

	clusterHostGRE := &network.ClusterHostGRE{
		BrName:  brName,
		GREKey:  greKey,
		Tunnels: hostGRETunnels,
	}

	return &provisionClusterHostGreWF{
		clustername:    clustername,
		ClusterHostGRE: clusterHostGRE,
		provisionAgent: pa,
		log:            wfLog,
	}
}

func (wf *provisionClusterHostGreWF) Run() error {
	wf.log.Infof("provisioning cluster %q host gre", wf.clustername)

	tasks, err := wf.getProvisionTask()
	if err != nil {
		wf.log.Errorf("error provisioning cluster %q host gre", wf.clustername)
		return err
	}

	for i := range tasks {
		go func(t task.Task) {
			wf.log.Infof("running provision cluster host gre task: %v", t)
			if err := t.Run(); err != nil {
				wf.log.Errorf("error running task: %s", err)
			}
		}(tasks[i])
	}
	return nil
}

func (wf *provisionClusterHostGreWF) getProvisionTask() ([]task.Task, error) {
	wf.log.Debugf("getting provision task for cluster %q host gre", wf.clustername)

	var tasks []task.Task
	tasks = wf.setupGRETasks(tasks)
	return tasks, nil
}

func (wf *provisionClusterHostGreWF) setupGRETasks(tasks []task.Task) []task.Task {
	wf.log.Debugf("setup gre tasks for cluster %q", wf.clustername)

	for _, tun := range wf.Tunnels {
		createTapDevTask := provisioner.NewCreateTapDev(tun.TapDevName, wf.log)
		createGREBridgeTask := provisioner.NewCreateGREBridge(wf.BrName, tun.TapDevName, wf.GREKey, tun.RemoteHostIP, wf.ovsClient, wf.log)

		greSeqTasks := task.NewSeqTask(wf.log, createTapDevTask, createGREBridgeTask)
		tasks = append(tasks, greSeqTasks)
	}
	return tasks
}
