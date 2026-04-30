# Manage example status page of the organization with ID 1.
resource "statuspal_status_page" "example" {
  organization_id = "1"
  status_page = {
    name      = "Example Terraform Status Page"
    url       = "example.com"
    time_zone = "Europe/Berlin"
  }
}

# Status page with a custom domain provisioned via Cloudflare.
# See https://github.com/statuspal/terraform-provider-statuspal/tree/main/examples/custom_domain
# for the full single-apply flow that wires domain_config together with the
# statuspal_domain_ssl_records and statuspal_custom_domain_validation waiter
# resources plus the corresponding cloudflare_record entries.
resource "statuspal_status_page" "with_custom_domain" {
  organization_id = "1"
  status_page = {
    name      = "Status Page with Custom Domain"
    url       = "example.com"
    time_zone = "UTC"

    domain_config = {
      provider = "cloudflare"
      domain   = "status.example.com"
    }
  }
}
