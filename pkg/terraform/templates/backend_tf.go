package templates

var Backend = `terraform {
  backend "s3" {
    bucket = "{{.Bucket}}"
    key    = "{{.Key}}"
    region = "{{.Region}}"
  }
}`
