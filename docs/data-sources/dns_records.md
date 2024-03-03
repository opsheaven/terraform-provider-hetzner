---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "hetzner_dns_records Data Source - terraform-provider-hetzner"
subcategory: ""
description: |-
  Hetzner Zone Records Data Source
---

# hetzner_dns_records (Data Source)

Hetzner Zone Records Data Source



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `zone_id` (String) Hetzner Zone Identifier

### Read-Only

- `records` (Attributes List) List of records created in the zone. (see [below for nested schema](#nestedatt--records))

<a id="nestedatt--records"></a>
### Nested Schema for `records`

Optional:

- `id` (String) Record Identifier
- `name` (String) Record name
- `zone_id` (String) Zone identifier that record belongs to

Read-Only:

- `ttl` (Number) Record TTL
- `type` (String) Record Type. Supported values: [ A,AAAA,NS,MX,CNAME,RP,TXT,SOA,HINFO,SRV,DANE,TLSA,DS,CAA ]
- `value` (String) Record value