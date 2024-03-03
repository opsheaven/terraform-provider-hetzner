package dns

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dsSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opsheaven/gohetznerdns"
)

type Zone struct {
	Id     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	NS     types.List   `tfsdk:"ns"`
	Paused types.Bool   `tfsdk:"paused"`
	Status types.String `tfsdk:"status"`
	TTL    types.Int64  `tfsdk:"ttl"`
}

type Zones struct {
	Name  types.String `tfsdk:"name"`
	Zones []Zone       `tfsdk:"zones"`
}

var ZoneDataSourceSchema = dsSchema.Schema{
	MarkdownDescription: "Hetzner Zone Data Source.",
	Attributes: map[string]dsSchema.Attribute{
		"id": dsSchema.StringAttribute{
			MarkdownDescription: "Zone Identifier",
			Optional:            true,
		},
		"name": dsSchema.StringAttribute{
			MarkdownDescription: "Zone Name",
			Optional:            true,
		},
		"ns": dsSchema.ListAttribute{
			MarkdownDescription: "Primary Nameservers assigned to the Zone. Managed by Hetzner.",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"paused": dsSchema.BoolAttribute{
			MarkdownDescription: "Zone activeness",
			Computed:            true,
		},
		"status": dsSchema.StringAttribute{
			MarkdownDescription: `
			Status of the zone. Supported values are:
			*verified*: Zone is verified.
			*failed*: Zone verification is failed.
			*pending*: Verification is in progress
			`,
			Computed: true,
		},
		"ttl": dsSchema.Int64Attribute{
			MarkdownDescription: "Zone Default TTL for zone records",
			Computed:            true,
		},
	},
}
var ZoneResourceSchema = rSchema.Schema{
	MarkdownDescription: "Hetzner Zone Resource.",
	Attributes: map[string]rSchema.Attribute{
		"id": rSchema.StringAttribute{
			MarkdownDescription: "Zone Identifier",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": rSchema.StringAttribute{
			MarkdownDescription: "Zone Name",
			Required:            true,
		},
		"ns": rSchema.ListAttribute{
			MarkdownDescription: "Primary Nameservers assigned to the Zone. Managed by Hetzner.",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"paused": rSchema.BoolAttribute{
			MarkdownDescription: "Zone activeness",
			Computed:            true,
		},
		"status": rSchema.StringAttribute{
			MarkdownDescription: `
			Status of the zone. Supported values are:
			*verified*: Zone is verified.
			*failed*: Zone verification is failed.
			*pending*: Verification is in progress
			`,
			Computed: true,
		},
		"ttl": rSchema.Int64Attribute{
			MarkdownDescription: "Zone Default TTL for zone records",
			Required:            true,
		},
	},
}
var ZonesDataSourceSchema = dsSchema.Schema{
	MarkdownDescription: "Hetzner Zones Data Source.",
	Attributes: map[string]dsSchema.Attribute{
		"name": dsSchema.StringAttribute{
			MarkdownDescription: "Zone full or partial name to query zones",
			Optional:            true,
		},
		"zones": dsSchema.ListNestedAttribute{
			MarkdownDescription: "List of Zones matches the given `1name ",
			Optional:            true,
			NestedObject: dsSchema.NestedAttributeObject{
				Attributes: ZoneDataSourceSchema.Attributes,
			},
		},
	},
}

func (z *Zone) mapFromHetznerZone(zone *gohetznerdns.Zone) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	if z.Id.IsNull() || z.Id.IsUnknown() {
		z.Id = types.StringValue(*zone.Id)
	}
	z.Name = types.StringValue(*zone.Name)
	z.Paused = types.BoolValue(*zone.Paused)
	z.Status = types.StringValue(*zone.Status)
	z.TTL = types.Int64Value(int64(*zone.TTL))

	elements := []attr.Value{}
	for _, ns := range zone.NS {
		elements = append(elements, types.StringValue(*ns))
	}
	z.NS, diagnostics = types.ListValue(types.StringType, elements)
	return diagnostics
}

func (z *Zones) mapFromHetznerZones(hetznerZones []*gohetznerdns.Zone) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}
	z.Zones = []Zone{}
	for _, hetznerZone := range hetznerZones {
		zone := Zone{}
		diagnostics.Append(zone.mapFromHetznerZone(hetznerZone)...)
		z.Zones = append(z.Zones, zone)
	}
	return diagnostics
}
