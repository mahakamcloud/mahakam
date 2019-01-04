package dns

// TFVars specifies the content of terraform.tfvars
var TFVars = `hostname = "[[.Name]]"
host = "[[.Host]]"
image_source_path = "[[.ImageSourcePath]]"
ssh_public_key = "[[.SSHPublicKey]]"
mac_address = "[[.MacAddress]]"
ip_address = "[[.IPAddress]]"
netmask = "[[.NetMask]]"
gateway = "[[.Gateway]]"
dns_address = "[[.DNSAddress]]"
dns_zone_name = "[[.DNSDomainName]]"
network_cidr = "[[.NetworkCIDR]]"
`
