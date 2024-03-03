package provider

import "github.com/opsheaven/gohetznerdns"

type dnsZones struct {
	Filter *string    `tfsdk:"filter"`
	Zones  []*dnsZone `tfsdk:"zones"`
}

type dnsZone struct {
	Id     *string   `tfsdk:"id"`
	Name   *string   `tfsdk:"name"`
	NS     []*string `tfsdk:"ns"`
	Status *string   `tfsdk:"status"`
	TTL    *int      `tfsdk:"ttl"`
}

type dnsZoneFile struct {
	Id       *string `tfsdk:"id"`
	Name     *string `tfsdk:"name"`
	ZoneFile *string `tfsdk:"zonefile"`
}

func newDnsZone(zone *gohetznerdns.Zone) *dnsZone {
	z := &dnsZone{}
	if zone != nil {
		z.mapFromZone(zone)
	}
	return z
}

func (zone *dnsZone) mapFromZone(z *gohetznerdns.Zone) {
	if z != nil {
		zone.Id = z.Id
		zone.Name = z.Name
		zone.NS = z.NS
		zone.Status = z.Status
		zone.TTL = z.TTL
	}
}
