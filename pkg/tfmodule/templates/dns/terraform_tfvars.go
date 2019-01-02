package dns

var TFVars = `hostname = "[[.Name]]"
host = "[[.Host]]"
image_source_path = "[[.ImageSourcePath]]"
ssh_public_key = "[[.SSHPublicKey]]"
mac_address = "[[.MacAddress]]"
ip_address = "[[.IPAddress]]"
netmask = "[[.NetMask]]"
gateway_ip = "[[.GatewayIP]]"
upstream_dns = "[[.UpstreamDNS]]"
dns_zone_name = "[[.DNSZoneName]]"
network_cidr = "${.NetworkCIDR}"
`
