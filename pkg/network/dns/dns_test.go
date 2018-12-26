package dns_test

import (
	"net"

	. "github.com/mahakamcloud/mahakam/pkg/network/dns"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DNS", func() {
	var (
		dnsConfig DNSConfig
	)

	BeforeEach(func() {
		dnsConfig = DNSConfig{
			PrivateIP:   net.ParseIP("10.10.10.1"),
			DNSZoneName: "sample.domain.io",
			Hostname:    "mahakam-dns-01",
		}
	})

	Describe("Generate zone file host record", func() {
		It("should return the correct zone file host record", func() {
			result := dnsConfig.GenerateZoneFileHostRecord()
			Expect(result).To(Equal("mahakam-dns-01 A 10.10.10.1"))
		})
	})

	Describe("Generate reverse zone file host record", func() {
		It("should return the correct reverse zone file host record", func() {
			result := dnsConfig.GenerateReverseZoneFileHostRecord()
			Expect(result).To(Equal("1 PTR mahakam-dns-01.sample.domain.io."))
		})
	})
})
