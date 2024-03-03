# Create new zone with TTL
resource "hetzner_dns_zone" "this" {
  name = "opsheaven.space"
  ttl  = 3600
}
