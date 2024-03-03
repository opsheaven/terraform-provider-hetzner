# Poulates all zones
data "hetzner_dns_zones" "all_zones" {
}

# Poulates all zones matching given keyword in the name field
data "hetzner_dns_zones" "zones_contains_opsheaven" {
  name = "opsheaven"
}
