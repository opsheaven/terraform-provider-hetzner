
# Get zone by name
data "hetzner_dns_zone" "this" {
  name = "opsheaven.space"
}

# Create new record
resource "hetzner_dns_record" "google_recovery_domain_verification" {
  name    = "@"
  type    = "TXT"
  ttl     = 3600
  value   = "google-gws-recovery-domain-verification=1111111"
  zone_id = hetzner_dns_zone.this.id
}
