package dns

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/opsheaven/gohetznerdns"
)

type ZoneService interface {
	List(zones *Zones) diag.Diagnostics
	Read(zone *Zone) diag.Diagnostics
	Create(zone *Zone) diag.Diagnostics
	Update(zone *Zone) diag.Diagnostics
	Delete(zone *Zone) diag.Diagnostics
}

type zoneServiceImpl struct {
	client gohetznerdns.ZoneService
}

var _ ZoneService = &zoneServiceImpl{}

func newZoneService(service gohetznerdns.ZoneService) ZoneService {
	return &zoneServiceImpl{client: service}
}

func (s *zoneServiceImpl) List(zones *Zones) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}
	var hetznerZones []*gohetznerdns.Zone
	var apiError error

	if !zones.Name.IsNull() {
		hetznerZones, apiError = s.client.GetAllZonesByName(zones.Name.ValueStringPointer())
	} else {
		hetznerZones, apiError = s.client.GetAllZones()
	}
	if apiError != nil {
		diagnostics.AddError("Hetzer Client Error", apiError.Error())
	} else {
		diagnostics.Append(zones.mapFromHetznerZones(hetznerZones)...)
	}
	return diagnostics
}

func (s *zoneServiceImpl) Read(zone *Zone) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}
	if !zone.Id.IsNull() && zone.Id.String() != "" {
		hetznerZone, err := s.client.GetZoneById(zone.Id.ValueStringPointer())
		if err != nil {
			diagnostics.AddError("Hetzer Client Error", err.Error())
		} else {
			zone.mapFromHetznerZone(hetznerZone)
		}
	} else if !zone.Name.IsNull() && zone.Name.String() != "" {
		hetznerZones, err := s.client.GetAllZonesByName(zone.Name.ValueStringPointer())
		if err != nil {
			diagnostics.AddError("Hetzer Client Error", err.Error())
		} else if len(hetznerZones) == 0 {
			diagnostics.AddError("Invalid Zone Name", fmt.Sprintf("Zone with %s can not be found", zone.Name.String()))
		} else if len(hetznerZones) > 1 {
			diagnostics.AddError("Multiple Zones", fmt.Sprintf("Found %d zones contains %s! Please provide full domain name", len(hetznerZones), zone.Name.String()))
		} else {
			diagnostics.Append(zone.mapFromHetznerZone(hetznerZones[0])...)
		}
	} else {
		diagnostics.AddError("Configuration Error", "ID or Name must be provided!")
	}

	return diagnostics
}

func (s *zoneServiceImpl) Create(zone *Zone) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}
	ttl := int(zone.TTL.ValueInt64())
	hetznerZone, err := s.client.CreateZone(&gohetznerdns.ZoneRequest{Name: zone.Name.ValueStringPointer(), TTL: &ttl})
	if err != nil {
		diagnostics.AddError("Hetzer Client Error", err.Error())
	} else {
		diagnostics.Append(zone.mapFromHetznerZone(hetznerZone)...)
	}
	return diagnostics
}

func (s *zoneServiceImpl) Update(zone *Zone) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}
	ttl := int(zone.TTL.ValueInt64())
	hetznerZone, err := s.client.UpdateZone(zone.Id.ValueStringPointer(), &gohetznerdns.ZoneRequest{Name: zone.Name.ValueStringPointer(), TTL: &ttl})
	if err != nil {
		diagnostics.AddError("Hetzer Client Error", err.Error())
	} else {
		diagnostics.Append(zone.mapFromHetznerZone(hetznerZone)...)
	}
	return diagnostics
}

func (s *zoneServiceImpl) Delete(zone *Zone) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}
	err := s.client.DeleteZone(zone.Id.ValueStringPointer())
	if err != nil {
		diagnostics.AddError("Hetzer Client Error", err.Error())
	}
	return diagnostics
}
