package dns

import (
	"net"
	"strings"

	"text/template"

	log "github.com/sirupsen/logrus"
)

// Config defines zone file configuration for a host
type Config struct {
	PrivateIP          net.IP
	DNSZoneName        string
	DNSReverseZonename string
	Hostname           string
}

const (
	zoneFileRecordTemplate        = `{{.Hostname}} A {{.PrivateIP}}`
	reverseZoneFileRecordTemplate = `{{$privateIP := .PrivateIP}}{{$arrPrivateIP := split $privateIP.String "."}}{{index $arrPrivateIP 3}} PTR {{.Hostname}}.{{.DNSZoneName}}.`
)

var (
	funcs = template.FuncMap{"split": strings.Split}
)

// GenerateZoneFileHostRecord returns a record to be appended in zone file
func (d Config) GenerateZoneFileHostRecord() string {
	dnsTemplVal := template.New("dns")
	tmpl, err := dnsTemplVal.Parse(zoneFileRecordTemplate)
	if err != nil {
		log.Errorf("Template parse error: %s", err)
	}

	var data strings.Builder
	err = tmpl.Execute(&data, d)
	if err != nil {
		log.Errorf("Template execute error: %s", err)
	}

	return data.String()
}

// GenerateReverseZoneFileHostRecord returns a record to be appended in reverse zone file
func (d Config) GenerateReverseZoneFileHostRecord() string {
	dnsTemplVal := template.New("dns").Funcs(funcs)
	tmpl, err := dnsTemplVal.Parse(reverseZoneFileRecordTemplate)
	if err != nil {
		log.Errorf("Template parse error: %s", err)
	}

	var data strings.Builder
	err = tmpl.Execute(&data, d)
	if err != nil {
		log.Errorf("Template execute error: %s", err)
	}

	return data.String()
}
