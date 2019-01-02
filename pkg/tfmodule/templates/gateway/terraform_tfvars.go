package gateway

var TFVars = `hostname = "[[.Name]]"
host = "[[.Host]]"
image_source_path = "[[.ImageSourcePath]]"
ssh_public_key = "[[.SSHPublicKey]]"
mac_address = "[[.MacAddress]]"
ip_address = "[[.IPAddress]]"
netmask = "[[.NetMask]]"
public_ip_address = "[[.PublicIPAddress]]"
public_netmask = "[[.PublicNetmask]]"
gateway = "[[.Gateway]]"
dns_address = "[[.DNSAddress]]"
dns_domain_name = "[[.DNSDomainName]]"
network_cidr = "[[.NetworkCIDR]]"
`
