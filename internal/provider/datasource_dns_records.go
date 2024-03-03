package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/opsheaven/terraform-provider-hetzner/internal/hetzner"
	"github.com/opsheaven/terraform-provider-hetzner/internal/hetzner/dns"
)

var (
	_ datasource.DataSource              = &dnsRecordsDataSource{}
	_ datasource.DataSourceWithConfigure = &dnsRecordsDataSource{}
)

type dnsRecordsDataSource struct {
	Service dns.RecordService
}

func NewRecordsDataSource() datasource.DataSource {
	return &dnsRecordsDataSource{}
}

func (datasource *dnsRecordsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
		d, diags := dataProvider.DNSServices()
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			datasource.Service = d.RecordService()
		}
	}
}

func (datasource *dnsRecordsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_records"
}

func (datasource *dnsRecordsDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dns.RecordsDataSourceSchema
}
func (datasource *dnsRecordsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dns.Records
	diags := req.Config.Get(ctx, &state)

	diags.Append(datasource.Service.List(&state)...)
	diags.Append(resp.State.Set(ctx, &state)...)
	resp.Diagnostics.Append(diags...)
}
