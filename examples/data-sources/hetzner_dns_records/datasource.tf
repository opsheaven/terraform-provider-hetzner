
# Get zone information by name
data "hetzner_dns_zone" "zone_by_name" {
  name = "opsheaven.space"
}

# Get all records owned by opsheaven.com
data "hetzner_dns_records" "opsheaven_all_records" {
  zone_id = "UFWX4H7TP93znuujDkzT9"
}
