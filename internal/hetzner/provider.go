package hetzner

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/opsheaven/terraform-provider-hetzner/internal/hetzner/dns"
)

type ProviderContext struct {
	DnsApiEnabled bool
	DnsApiToken   string
}

type Provider interface {
	DNSServices() (dns.DNSServices, diag.Diagnostics)
}

type provider struct {
	context     *ProviderContext
	dnsServices dns.DNSServices
}

var _ Provider = &provider{}

func NewProvider(ctx *ProviderContext) (Provider, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	provider := &provider{context: ctx}

	if ctx.DnsApiEnabled {
		dns, diags := dns.NewClient(ctx.DnsApiToken)
		diagnostics.Append(diags...)
		if !diagnostics.HasError() {
			provider.dnsServices = dns
		} else {
			diagnostics.AddWarning("DNS Api Disabled", "DNS Service has initialization errors.")
		}
	} else {
		diagnostics.AddWarning("DNS Api Disabled", "DNS Service is disabled and all DNS calls will generate error")
	}
	return provider, diagnostics
}

func (p *provider) DNSServices() (dns.DNSServices, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	if p.context.DnsApiEnabled {
		return p.dnsServices, diagnostics
	}
	diagnostics.AddError("DNS Api Disabled", "DNS Api is disabled, please enable on the provider")
	return nil, diagnostics
}
