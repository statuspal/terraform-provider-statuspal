# Manage example service of the status page with subdomain "example-com".
resource "statuspal_service" "example" {
  status_page_subdomain = "example-com"
  service = {
    name = "Example Terraform Service"
  }
}
