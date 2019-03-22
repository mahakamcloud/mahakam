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
	MahakamAPIBasePath = "/v1"

	// Custom path for storing resource in kvstore
	KeyPathMahakam       = "mahakamcloud/"
	KeyPathNetworkSubnet = KeyPathMahakam + "network/subnets/"
	KeyPathNetworkIPPool = KeyPathMahakam + "network/ip-pools/"

	// Helm default configuration
	HelmDefaultNamespace             = "default"
	HelmDefaultTillerNamespace       = "kube-system"
	HelmDefaultKubecontext           = "kubernetes-admin@kubernetes"
	HelmControllerWait               = false
	HelmControllerDefaultWaitTimeout = 300
	HelmDefaultChartValuesDirectory  = "/opt/mahakamcloud/chartvalues/"

	// Default mahakam config
	MahakamMultiKubeconfigPath     = "/opt/mahakamcloud/clusters"
	MahakamSSHPrivateKeyPath       = "/root/.ssh/id_rsa"
	MahakamDefaultNetworkInterface = "ens3"

	// Default terraform config
	TerraformDefaultDirectory = "/opt/mahakamcloud/terraform/"
	TerraformDefaultBucket    = "tf-mahakam"
	TerraformDefaultRegion    = "ap-southeast-1"

	// Default kubernetes node config
	KubernetesNodeUsername        = "ubuntu"
	KubernetesAdminKubeconfigPath = "/home/ubuntu/.kube/config"
	KubernetesAPIServerPort       = 6443

	// Node Health check config
	NodePingTimeout = 5 * time.Second
	NodePingRetry   = 20
	NodePingDelay   = 60 * time.Second

	// Storage Backend Health check config
	StorageBackendPingTimeout = 3 * time.Second
	StorageBackendPingRetry   = 5
	StorageBackendPingDelay   = 2 * time.Second

	// Various keys for storing metadata in map
	KeyControlPlaneIP     = "key-control-plane-ip"
	KeyPodNetworkCidr     = "key-pod-network-cidr"
	KeyKubeadmToken       = "key-kubeadm-token"
	KeyClusterNetworkCidr = "key-cluster-network-cidr"
	KeySubnetMask         = "key-subnet-mask"
	KeySubnetAddress      = "key-subnet-address"
	KeyBroadcastAddress   = "key-broadcast-address"

	// Enum for ip pool service
	IPPoolActionAllocate = "ALLOCATE"
	IPPoolActionRelease  = "RELEASE"

	TerraformTFState = "/terraform.tfstate"
)
