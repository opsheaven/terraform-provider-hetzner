package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opsheaven/terraform-provider-hetzner/internal/hetzner"
)

var (
	_ provider.Provider = &hetznerProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &hetznerProvider{
			typeName: "hetzner",
			version:  version,
		}
	}
}

type hetznerProvider struct {
	typeName string
	version  string
}
type hetznerProviderModel struct {
	DnsApiEnabled types.Bool   `tfsdk:"dns_api_enabled"`
	DnsApiToken   types.String `tfsdk:"dns_api_token"`
}

func (p *hetznerProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = p.typeName
	resp.Version = p.version
}

func (p *hetznerProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `This providers helps you to manage of Hetzner resources`,
		Attributes: map[string]schema.Attribute{
			"dns_api_enabled": schema.BoolAttribute{
				MarkdownDescription: "Optional dns api enabler flag. DNS API is enabled by default.",
				Optional:            true,
			},
			"dns_api_token": schema.StringAttribute{
				MarkdownDescription: "Optional DNS Api Authentication token. When missing provider will populate it from `HETZNER_DNS_API_TOKEN` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *hetznerProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config hetznerProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.DnsApiEnabled.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("dns_api_enabled"),
			"Hetzner DNS API Enabler Flag",
			"The provider cannot create the Hezner DNS API client as there is an unknown configuration value for the HashiCups DNS API.",
		)
	}

	if config.DnsApiToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Hetzner DNS API Token",
			"The provider cannot create the Hezner DNS API client as there is an unknown configuration value for the Hetzner DNS API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HETZNER_DNS_API_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	dns_api_enabled := true
	dns_api_token := os.Getenv("HETZNER_DNS_API_TOKEN")

	if !config.DnsApiEnabled.IsNull() {
		dns_api_enabled = config.DnsApiEnabled.ValueBool()
	}
	if !config.DnsApiToken.IsNull() {
		dns_api_token = config.DnsApiToken.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if dns_api_enabled && dns_api_token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("dns_api_token"),
			"Missing Hetzner DNS API Token",
			"The provider cannot create the Hetzner DNS API client as there is a missing or empty value for the Hetzner DNS API token. "+
				"Set the host value in the configuration or use the HETZNER_DNS_API_TOKEN environment variable. "+
				"To disable DNS Api lookups, disable DNS API by configuring `dns_api_enabled=false`",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	provider, providerDiags := hetzner.NewProvider(&hetzner.ProviderContext{DnsApiEnabled: dns_api_enabled, DnsApiToken: dns_api_token})
	resp.Diagnostics.Append(providerDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.DataSourceData = provider
	resp.ResourceData = provider
}

func (p *hetznerProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewZonesDataSource,
		NewZoneDataSource,
		NewRecordsDataSource,
	}
}

func (p *hetznerProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDnsZoneResource,
		NewDnsRecordResource,
	}
}
