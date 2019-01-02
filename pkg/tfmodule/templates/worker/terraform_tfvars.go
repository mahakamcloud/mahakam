package worker

var TFVars = `hostname = "[[.Name]]"
host = "[[.Host]]"
image_source_path = "[[.ImageSourcePath]]"
ssh_public_key = "[[.SSHPublicKey]]"
mac_address = "[[.MacAddress]]"
ip_address = "[[.IPAddress]]"
netmask = "[[.NetMask]]"
gateway = "[[.Gateway]]"
dns_address = "[[.DNSAddress]]"
dns_domain_name = "[[.DNSDomainName]]"
controlplane_ip = "[[.ControlPlaneIP]]"
kubeadm_token = "[[.KubeadmToken]]"
`
