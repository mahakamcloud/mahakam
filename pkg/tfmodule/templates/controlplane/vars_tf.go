package controlplane

var Vars = `variable "hostname" {
    type = "string"
}

variable "host" {
    type = "string"
}

variable "image_source_path" {
    type = "string"
}

variable "password" {
    type    = "string"
    default = "passw0rd"
}

variable "mac_address" {
    type = "string"
}

variable "ip_address" {
    type = "string"
}

variable "netmask" {
    type = "string"
}

variable "gateway" {
    type = "string"
}

variable "dns_address" {
    type = "string"
}

variable "dns_domain_name" {
    type = "string"
}

variable "pod_network_cidr" {
    type = "string"
}

variable "kubeadm_token" {
    type = "string"
}
`
