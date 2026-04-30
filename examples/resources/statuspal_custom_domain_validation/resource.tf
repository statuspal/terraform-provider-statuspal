# Waiter resource that polls until a status page's custom domain reaches
# "active" status. Use this as the final step in a custom domain setup to
# confirm the domain is fully configured before proceeding.

resource "statuspal_custom_domain_validation" "example" {
  organization_id       = "1"
  status_page_subdomain = "example-subdomain"
  timeout_seconds       = 600
}
