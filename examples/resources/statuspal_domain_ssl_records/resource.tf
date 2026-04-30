# Waiter resource that polls until the SSL certificate DNS challenge records
# are available for a status page with a custom domain configured via Cloudflare.
#
# Use this between creating the CNAME routing record and the TXT certificate
# record to complete custom domain setup in a single terraform apply.

resource "statuspal_domain_ssl_records" "example" {
  organization_id       = "1"
  status_page_subdomain = "example-subdomain"
  timeout_seconds       = 300
}
