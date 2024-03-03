package dns

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/opsheaven/gohetznerdns"
)

type RecordService interface {
	List(records *Records) diag.Diagnostics
	Read(record *Record) diag.Diagnostics
	Create(record *Record) diag.Diagnostics
	Update(record *Record) diag.Diagnostics
	Delete(record *Record) diag.Diagnostics
}

type recordServiceImpl struct {
	client gohetznerdns.RecordService
}

var _ RecordService = &recordServiceImpl{}

func newRecordService(service gohetznerdns.RecordService) RecordService {
	return &recordServiceImpl{client: service}
}

func (s *recordServiceImpl) List(records *Records) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}
	hetznerRecords, err := s.client.GetAllRecords(records.ZoneId.ValueStringPointer())

	if err != nil {
		diagnostics.AddError("Hetzer Client Error", err.Error())
	} else {
		diagnostics.Append(records.mapFromHetznerRecords(hetznerRecords)...)
	}

	return diagnostics
}

func (s *recordServiceImpl) Read(record *Record) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}
	if !record.Id.IsNull() && record.Id.String() != "" {
		hetznerRecord, err := s.client.GetRecord(record.Id.ValueStringPointer())
		if err != nil {
			diagnostics.AddError("Hetzer Client Error", err.Error())
		} else {
			diagnostics.Append(record.mapFromHetznerRecord(hetznerRecord)...)
		}
	} else {
		diagnostics.AddError("Configuration Error", "ID must be provided!")
	}
	return diagnostics
}

func (s *recordServiceImpl) Create(record *Record) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}
	ttl := int(record.TTL.ValueInt64())
	hetznerRecord := &gohetznerdns.Record{
		Type:   record.Type.ValueStringPointer(),
		ZoneId: record.ZoneId.ValueStringPointer(),
		Name:   record.Name.ValueStringPointer(),
		TTL:    &ttl,
	}
	var value string
	if record.Type.String() == "TXT" {
		value = fmt.Sprintf("\"%s\"", record.Value.ValueString())
	} else {
		value = record.Value.ValueString()
	}
	hetznerRecord.Value = &value

	hetznerRecord, err := s.client.CreateRecord(hetznerRecord)
	if err != nil {
		diagnostics.AddError("Hetzer Client Error", err.Error())
	} else {
		diagnostics.Append(record.mapFromHetznerRecord(hetznerRecord)...)
	}
	return diagnostics
}

func (s *recordServiceImpl) Update(record *Record) diag.Diagnostics {
	ttl := int(record.TTL.ValueInt64())
	diagnostics := diag.Diagnostics{}
	hetznerRecord := &gohetznerdns.Record{
		Id:     record.Id.ValueStringPointer(),
		Type:   record.Type.ValueStringPointer(),
		ZoneId: record.ZoneId.ValueStringPointer(),
		Name:   record.Name.ValueStringPointer(),
		TTL:    &ttl,
	}
	var value string
	if record.Type.String() == "TXT" {
		value = fmt.Sprintf("\"%s\"", record.Value.ValueString())
	} else {
		value = record.Value.ValueString()
	}
	hetznerRecord.Value = &value
	hetznerRecord, err := s.client.UpdateRecord(hetznerRecord)
	if err != nil {
		diagnostics.AddError("Hetzer Client Error", err.Error())
	} else {
		if hetznerRecord.TTL == nil {

			hetznerRecord.TTL = &ttl
		}
		diagnostics.Append(record.mapFromHetznerRecord(hetznerRecord)...)
	}
	return diagnostics
}

func (s *recordServiceImpl) Delete(record *Record) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}
	err := s.client.DeleteRecord(record.Id.ValueStringPointer())
	if err != nil {
		diagnostics.AddError("Hetzer Client Error", err.Error())
	}
	return diagnostics
}
