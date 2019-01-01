package gateway

var TFVars = `hostname = "{{.Name}}"
host = "{{.Host}}"
image_source_path = "{{.ImageSourcePath}}"
ssh_public_key = "{{.SSHPublicKey}}"
mac_address = "{{.MacAddress}}"
ip_address = "{{.IPAddress}}"
netmask = "{{.NetMask}}"
public_mac_address = "{{.PublicMacAddress}}"
public_ip_address = "{{.PublicIPAddress}}"
public_netmask = "{{.PublicNetmask}}"
gateway = "{{.Gateway}}"
dns_address = "{{.DNSAddress}}"
dns_zone_name = "{{.DNSDomainName}}"
network_cidr = "${.NetworkCIDR}"
`
