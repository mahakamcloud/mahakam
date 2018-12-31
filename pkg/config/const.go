package config

import "time"

// Some config constants or environment variables will go away
// once we populate this through kind of config.yaml
const (
	// ResourceOwner hardcodes all tenant resources to be owned by gojek
	// since we don't have auth mechanism yet
	ResourceOwnerGojek   = "gojek"
	ResourceOwnerMahakam = "mahakam"

	// Mahakam API server config
	MahakamAPIBasePath    = "/v1"
	MahakamAPIDefaultPort = 9000

	// Custom path for storing resource in kvstore
	KeyPathMahakam       = "mahakamcloud/"
	KeyPathNetworkSubnet = KeyPathMahakam + "network/subnets/"

	// Helm default configuration
	HelmDefaultNamespace             = "default"
	HelmDefaultTillerNamespace       = "kube-system"
	HelmDefaultKubecontext           = "kubernetes-admin@kubernetes"
	HelmControllerWait               = false
	HelmControllerDefaultWaitTimeout = 300
	HelmDefaultChartValuesDirectory  = "/opt/mahakamcloud/chartvalues/"

	// Default mahakam config path to store multiple kubeconfig files
	MahakamMultiKubeconfigPath = "/opt/mahakamcloud/clusters"
	MahakamSSHPrivateKeyPath   = "/root/.ssh/id_rsa"

	// Default terraform config
	TerraformDefaultDirectory = "/opt/mahakamcloud/terraform/"
	TerraformDefaultBucket    = "tf-mahakam"
	TerraformDefaultRegion    = "ap-southeast-1"

	// Default kubernetes node config
	KubernetesNodeUsername        = "ubuntu"
	KubernetesAdminKubeconfigPath = "/home/ubuntu/.kube/config"
	KubernetesAPIServerPort       = 6443
	KubernetesNodePingTimeout     = 5 * time.Second
	KubernetesNodePingRetry       = 20
	KubernetesNodePingDelay       = 60 * time.Second

	// Various keys for storing metadata in map
	KeyControlPlaneIP = "key-control-plane-ip"
	KeyPodNetworkCidr = "key-pod-network-cidr"
	KeyKubeadmToken   = "key-kubeadm-token"
)
