# Get zone information by id
data "hetzner_dns_zone" "zone_by_id" {
  id = "UFWX4H7TP93znuujDkzT9"
}

# Get zone by name
data "hetzner_dns_zone" "zone_by_name" {
  name = "opsheaven.space"
}
