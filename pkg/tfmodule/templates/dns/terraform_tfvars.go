package dns

var TFVars = `hostname = "{{.Name}}"
host = "{{.Host}}"
image_source_path = "{{.ImageSourcePath}}"
ssh_public_key = "{{.SSHPublicKey}}"
mac_address = "{{.MacAddress}}"
primary_ip_address = "{{.PrimaryIPAddress}}"
replica_ip_address = "{{.ReplicaIPAddress}}"
netmask = "{{.NetMask}}"
gateway = "{{.Gateway}}"
dns_address = "{{.DNSAddress}}"
dns_zone_name = "{{.DNSDomainName}}"
network_cidr = "${.NetworkCIDR}"
`
