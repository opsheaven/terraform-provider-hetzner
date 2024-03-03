package dns

import (
	"context"
	"fmt"
	"slices"
	"strings"

	dsSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opsheaven/gohetznerdns"
)

type Record struct {
	Id     types.String `tfsdk:"id"`
	Type   types.String `tfsdk:"type"`
	ZoneId types.String `tfsdk:"zone_id"`
	Name   types.String `tfsdk:"name"`
	Value  types.String `tfsdk:"value"`
	TTL    types.Int64  `tfsdk:"ttl"`
}

type Records struct {
	ZoneId  types.String `tfsdk:"zone_id"`
	Records []Record     `tfsdk:"records"`
}

var RecordDataSourceSchema = dsSchema.Schema{
	MarkdownDescription: "Hetzner Record DataSource",
	Attributes: map[string]dsSchema.Attribute{
		"id": dsSchema.StringAttribute{
			MarkdownDescription: "Record Identifier",
			Optional:            true,
			Computed:            true,
		},
		"type": dsSchema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("Record Type. Supported values: [ %s ]", strings.Join(allowedRecordTypes, ",")),
			Computed:            true,
		},
		"zone_id": dsSchema.StringAttribute{
			MarkdownDescription: "Zone identifier that record belongs to",
			Optional:            true,
			Computed:            true,
		},
		"name": dsSchema.StringAttribute{
			MarkdownDescription: "Record name",
			Computed:            true,
			Optional:            true,
		},
		"value": dsSchema.StringAttribute{
			MarkdownDescription: "Record value",
			Computed:            true,
		},
		"ttl": dsSchema.Int64Attribute{
			MarkdownDescription: "Record TTL",
			Computed:            true,
		},
	},
}

var RecordResourceSchema = rSchema.Schema{
	MarkdownDescription: "Hetzner Record Resource.",
	Attributes: map[string]rSchema.Attribute{
		"id": rSchema.StringAttribute{
			MarkdownDescription: "Record Identifier",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"type": rSchema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("Record Type. Supported values: [ %s ]", strings.Join(allowedRecordTypes, ",")),
			Required:            true,
			Validators: []validator.String{
				&recordTypeValidator{},
			},
		},
		"zone_id": rSchema.StringAttribute{
			MarkdownDescription: "Zone identifier that record belongs to",
			Required:            true,
		},
		"name": rSchema.StringAttribute{
			MarkdownDescription: "Record name",
			Required:            true,
		},
		"value": rSchema.StringAttribute{
			MarkdownDescription: "Record value",
			Required:            true,
		},
		"ttl": rSchema.Int64Attribute{
			MarkdownDescription: "Record TTL",
			Required:            true,
		},
	},
}

var RecordsDataSourceSchema = dsSchema.Schema{
	MarkdownDescription: "Hetzner Zone Records Data Source",
	Attributes: map[string]dsSchema.Attribute{
		"zone_id": dsSchema.StringAttribute{
			MarkdownDescription: "Hetzner Zone Identifier",
			Required:            true,
		},
		"records": dsSchema.ListNestedAttribute{
			MarkdownDescription: "List of records created in the zone.",
			Computed:            true,
			NestedObject: dsSchema.NestedAttributeObject{
				Attributes: RecordDataSourceSchema.Attributes,
			},
		},
	},
}

func (r *Record) mapFromHetznerRecord(record *gohetznerdns.Record) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	r.Id = types.StringValue(*record.Id)
	r.Name = types.StringValue(*record.Name)
	if record.TTL != nil {
		r.TTL = types.Int64Value(int64(*record.TTL))
	}
	r.Type = types.StringValue(*record.Type)
	if r.Type.String() == "TXT" {
		r.Value = types.StringValue(*record.Value)
	} else {
		r.Value = types.StringValue(strings.Trim(*record.Value, "\""))
	}
	r.ZoneId = types.StringValue(*record.ZoneId)
	return diagnostics
}

func (r *Records) mapFromHetznerRecords(hetznerRecords []*gohetznerdns.Record) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}
	r.Records = []Record{}
	for _, hetznerRecord := range hetznerRecords {
		record := Record{}
		diagnostics.Append(record.mapFromHetznerRecord(hetznerRecord)...)
		r.Records = append(r.Records, record)
	}
	return diagnostics
}

var allowedRecordTypes = []string{"A", "AAAA", "NS", "MX", "CNAME", "RP", "TXT", "SOA", "HINFO", "SRV", "DANE", "TLSA", "DS", "CAA"}

type recordTypeValidator struct {
}

func (r *recordTypeValidator) Description(context.Context) string {
	return "Type of the Type field"
}

func (r *recordTypeValidator) MarkdownDescription(context.Context) string {
	return fmt.Sprintf("Must be one of the values - [%s]", strings.Join(allowedRecordTypes, ","))
}

func (r *recordTypeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		resp.Diagnostics.AddError(fmt.Sprintf("%s is unknown", req.Path.String()), fmt.Sprintf("Please provide supported value for the field. Supported values: [%s]", strings.Join(allowedRecordTypes, ",")))
	} else if !slices.Contains(allowedRecordTypes, req.ConfigValue.ValueString()) {
		resp.Diagnostics.AddError(fmt.Sprintf("Invalid Value for %s", req.Path.String()), fmt.Sprintf("Please provide supported value for the field. Current Value: %s, Supported values: [%s]", req.ConfigValue.ValueString(), strings.Join(allowedRecordTypes, ",")))
	}
}
