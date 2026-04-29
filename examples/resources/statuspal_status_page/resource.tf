# Basic status page with no custom domain.
resource "statuspal_status_page" "example" {
  organization_id = "1"
  status_page = {
    name      = "Example Terraform Status Page"
    url       = "example.com"
    time_zone = "Europe/Berlin"
  }
}

# Status page with a Cloudflare-backed custom domain. After apply, read the
# computed `domain_config.validation_records` to find the DNS records you need
# to create on your DNS provider so Cloudflare can validate ownership and
# issue a certificate.
resource "statuspal_status_page" "with_cloudflare_domain" {
  organization_id = "1"
  status_page = {
    name      = "Example with Cloudflare Custom Domain"
    url       = "example.com"
    time_zone = "Europe/Berlin"

    domain_config = {
      provider = "cloudflare"
      domain   = "status.example.com"
    }
  }
}

# Status page with a Bunny CDN-backed custom domain. The required DNS record
# (a single CNAME) is exposed via `domain_config.validation_records`.
resource "statuspal_status_page" "with_bunny_domain" {
  organization_id = "1"
  status_page = {
    name      = "Example with Bunny Custom Domain"
    url       = "example.com"
    time_zone = "Europe/Berlin"

    domain_config = {
      provider = "bunny"
      domain   = "status.example.com"
    }
  }
}
