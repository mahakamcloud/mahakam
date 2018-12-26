package dns

import (
	"log"
	"net"
	"strings"

	"text/template"
)

type DNSConfig struct {
	PrivateIP          net.IP
	DNSZoneName        string
	DNSReverseZonename string
	Hostname           string
}

const (
	zoneFileRecordTemplate = `{{.Hostname}} A {{.PrivateIP}}`
)

// GenerateZoneFileHostRecord returns a record to be appended in zone file
func (d DNSConfig) GenerateZoneFileHostRecord() string {
	dnsTemplVal := template.New("dns")
	tmpl, err := dnsTemplVal.Parse(zoneFileRecordTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
	}

	var data strings.Builder
	err = tmpl.Execute(&data, d)
	if err != nil {
		log.Fatal("Execute: ", err)
	}

	return data.String()
}
