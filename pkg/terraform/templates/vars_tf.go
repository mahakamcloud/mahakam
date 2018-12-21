package templates

var Vars = `variable "dns_dhcp_server_username" {
    type = "string"
}

variable "hostname" {
    type = "string"
}

variable "password" {
    type    = "string"
    default = "passw0rd"
}

variable "ip_address" {
    type = "string"
}

variable "netmask" {
    type = "string"
}

variable "host" {
    type = "string"
}

variable "gate_nss_api_key" {
    type = "string"
}

variable "image_source_path" {
    type = "string"
}

variable "mac_address" {
    type = "string"
}

variable "dns_address" {
    type = "string"
}`
