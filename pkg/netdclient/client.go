package netdclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/netd/netd/host"
)

const (
	NetdHostBaseUrl      = "http://%s:%s/v1"
	NetdPort             = "80"
	NetdCreateNetworkAPI = "/network"
	contentTypeJSON      = "application/json"
)

type Network struct {
	Name string `json:"name"`
	Key  int    `json:"key"`
}

type Host struct {
	Name       string `json:"name"`
	IPAddr     string `json:"ip"`
	IPMaskSize string `json:"ipMask"`
}

type CreateNetworkRequest struct {
	Network Network `json:"network"`
	Hosts   []*Host `json:"hosts"`
}

type bridgeResp struct {
	Name   string `json:"name"`
	Status bool   `json:"status"`
	Err    string `json:"error"`
}

type greTunnelResp struct {
	Name   string     `json:"name"`
	Host   *host.Host `json:"host"`
	Status bool       `json:"status"`
	Err    string     `json:"error"`
}

type libvirtNetResp struct {
	Name   string `json:"name"`
	Status bool   `json:"status"`
	Err    string `json:"error"`
}

type CreateNetworkResponse struct {
	Status         bool             `json:"status"`
	BridgeResp     *bridgeResp      `json:"bridge"`
	GRETunnelsResp []*greTunnelResp `json:"gre_tunnels"`
	LibvirtNetResp *libvirtNetResp  `json:"libvirtnet"`
}

type Client struct{}

// CreateNetwork creates network on provided BareMetalHosts
func (s *Client) CreateNetwork(hosts []*models.BareMetalHost, network *models.GreNetwork) [][]byte {
	var netdResponses [][]byte

	// TODO(vjdhama) : make concurrent calls to netd hosts
	for _, h := range hosts {
		netdPayload, err := constructPayload(remove(hosts, h), network)
		if err == nil {
			netdResponses = append(netdResponses, netdPayload)
		}

		net, err := s.createNetwork(netdPayload, string(*h.IP), NetdPort)
		if err == nil {
			netdResponses = append(netdResponses, net)
		}
	}

	return netdResponses
}

func constructPayload(h []*models.BareMetalHost, n *models.GreNetwork) ([]byte, error) {
	var hosts []*Host

	for _, host := range h {
		netdHost := &Host{
			Name:       string(*host.Name),
			IPAddr:     string(*host.IP),
			IPMaskSize: string(*host.IPMask),
		}

		hosts = append(hosts, netdHost)
	}

	network := &Network{
		Name: string(*n.Name),
		Key:  int(n.GREKey),
	}

	request := &CreateNetworkRequest{
		Network: *network,
		Hosts:   hosts,
	}

	return json.Marshal(request)
}

func (s *Client) createNetwork(reader []byte, hostIP, hostPort string) ([]byte, error) {
	client := &http.Client{}

	url := fmt.Sprintf(NetdHostBaseUrl+NetdCreateNetworkAPI, hostIP, hostPort)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reader))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentTypeJSON)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func remove(hosts []*models.BareMetalHost, host *models.BareMetalHost) []*models.BareMetalHost {
	for i, v := range hosts {
		if v == host {
			hosts[i] = hosts[len(hosts)-1]
			return hosts[:len(hosts)-1]
		}
	}

	return hosts
}
