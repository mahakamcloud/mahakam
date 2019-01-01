package dhcp

var TFVars = `hostname = "{{.Name}}"
host = "{{.Host}}"
image_source_path = "{{.ImageSourcePath}}"
ssh_public_key = "{{.SSHPublicKey}}"
mac_address = "{{.MacAddress}}"
ip_address = "{{.IPAddress}}"
netmask = "{{.NetMask}}"
gateway = "{{.Gateway}}"
dns_address = "{{.DNSAddress}}"
dns_zone_name = "{{.DNSZoneName}}"
network_cidr = "${.NetworkCIDR}"
broadcast_address = "${.BroadcastAddress}"
subnet_address = "${.SubnetAddress}"
subnet_mask = "${.SubnetMask}"
`
