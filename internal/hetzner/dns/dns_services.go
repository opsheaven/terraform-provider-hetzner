package dns

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/opsheaven/gohetznerdns"
)

type DNSServices interface {
	ZoneService() ZoneService
	RecordService() RecordService
}

type dnsServicesImpl struct {
	recordService RecordService
	zoneService   ZoneService
}

var _ DNSServices = &dnsServicesImpl{}

func (d *dnsServicesImpl) RecordService() RecordService {
	return d.recordService
}

func (d *dnsServicesImpl) ZoneService() ZoneService {
	return d.zoneService
}

func NewClient(dnsApiToken string) (DNSServices, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}

	dnsClient, err := gohetznerdns.NewClient(dnsApiToken)
	if err != nil {
		diagnostics.AddError("DNS Client Initialization Error", err.Error())
	}
	if diagnostics.HasError() {
		return nil, diagnostics
	}
	return &dnsServicesImpl{
		recordService: newRecordService(dnsClient.GetRecordService()),
		zoneService:   newZoneService(dnsClient.GetZoneService()),
	}, diagnostics
}
