package gateway

// Vars specifies the terraform variables template
var Vars = `variable "hostname" {
    type = "string"
}

variable "libvirt_host" {
    type = "string"
}

variable "image_source_path" {
    type = "string"
}

variable "password" {
    type    = "string"
    default = "passw0rd"
}

variable "ssh_public_key" {
    type = "string"
}

variable "mac_address" {
    type = "string"
}

variable "ip_address" {
    type = "string"
}

variable "public_ip_address" {
    type = "string"
}

variable "netmask" {
    type = "string"
}

variable "public_netmask" {
    type = "string"
}

variable "gateway_ip" {
    type = "string"
}

variable "public_gateway_ip" {
    type = "string"
}

variable "dns_server" {
    type = "string"
}

variable "dns_zone_name" {
    type = "string"
}

variable "network_cidr" {
    type = "string"
}
`
