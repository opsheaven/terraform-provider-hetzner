package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/opsheaven/terraform-provider-hetzner/internal/hetzner"
	"github.com/opsheaven/terraform-provider-hetzner/internal/hetzner/dns"
)

var (
	_ datasource.DataSource              = &dnsZoneDataSource{}
	_ datasource.DataSourceWithConfigure = &dnsZoneDataSource{}
)

type dnsZoneDataSource struct {
	Service dns.ZoneService
}

func NewZoneDataSource() datasource.DataSource {
	return &dnsZoneDataSource{}
}

func (datasource *dnsZoneDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	dataProvider, ok := req.ProviderData.(hetzner.Provider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hetznerDataProvider, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	} else {
		service, diags := dataProvider.DNSServices()
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			datasource.Service = service.ZoneService()
		}
	}
}

func (datasource *dnsZoneDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_zone"
}

func (d *dnsZoneDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dns.ZoneDataSourceSchema
}

func (datasource *dnsZoneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dns.Zone
	diags := req.Config.Get(ctx, &state)
	diags.Append(datasource.Service.Read(&state)...)
	diags.Append(resp.State.Set(ctx, &state)...)
	resp.Diagnostics.Append(diags...)
}
