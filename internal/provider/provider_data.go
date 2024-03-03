package provider

import (
	"github.com/opsheaven/gohetznerdns"
)

type hetznerDataProvider struct {
	dnsApiEnabled bool
	dnsClient     gohetznerdns.HetznerDNS
}

func newHetznerDataProvider(dns_api_enabled bool, dns_api_token string) (*hetznerDataProvider, error) {
	dataProvider := &hetznerDataProvider{
		dnsApiEnabled: false,
		dnsClient:     nil,
	}

	if dns_api_enabled {
		dataProvider.dnsApiEnabled = true
		dnsClient, err := gohetznerdns.NewClient(dns_api_token)
		if err != nil {
			return nil, err
		}
		dataProvider.dnsClient = dnsClient
	}
	return dataProvider, nil
}
