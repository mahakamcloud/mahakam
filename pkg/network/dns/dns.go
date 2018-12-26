package dns

import (
	"net"
	"strings"

	"text/template"

	log "github.com/sirupsen/logrus"
)

type DNSConfig struct {
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
func (d DNSConfig) GenerateZoneFileHostRecord() string {
	dnsTemplVal := template.New("dns")
	tmpl, err := dnsTemplVal.Parse(zoneFileRecordTemplate)
	if err != nil {
		log.Errorf("Template parse error: ", err)
	}

	var data strings.Builder
	err = tmpl.Execute(&data, d)
	if err != nil {
		log.Errorf("Template execute error: ", err)
	}

	return data.String()
}

// GenerateReverseZoneFileHostRecord returns a record to be appended in reverse zone file
func (d DNSConfig) GenerateReverseZoneFileHostRecord() string {
	dnsTemplVal := template.New("dns").Funcs(funcs)
	tmpl, err := dnsTemplVal.Parse(reverseZoneFileRecordTemplate)
	if err != nil {
		log.Errorf("Parse: ", err)
	}

	var data strings.Builder
	err = tmpl.Execute(&data, d)
	if err != nil {
		log.Errorf("Execute: ", err)
	}

	return data.String()
}
