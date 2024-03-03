package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/opsheaven/terraform-provider-hetzner/internal/hetzner"
	"github.com/opsheaven/terraform-provider-hetzner/internal/hetzner/dns"
)

var _ resource.Resource = &dnsZoneResource{}
var _ resource.ResourceWithConfigure = &dnsZoneResource{}
var _ resource.ResourceWithImportState = &dnsZoneResource{}

type dnsZoneResource struct {
	Service dns.ZoneService
}

func NewDnsZoneResource() resource.Resource {
	return &dnsZoneResource{}
}

func (resource *dnsZoneResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	dataProvider, ok := req.ProviderData.(hetzner.Provider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *hetznerDataProvider, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	} else {
		service, diags := dataProvider.DNSServices()
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			resource.Service = service.ZoneService()
		}
	}
}
func (resource *dnsZoneResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_zone"
}

func (resource *dnsZoneResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = dns.ZoneResourceSchema
}

func (resource *dnsZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state dns.Zone
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	resp.Diagnostics.Append(resource.Service.Create(&state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (resource *dnsZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dns.Zone
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resource.Service.Read(&state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (resource *dnsZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state dns.Zone
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	resp.Diagnostics.Append(resource.Service.Update(&state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (resource *dnsZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state dns.Zone
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(resource.Service.Delete(&state)...)
}

func (r *dnsZoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
