# Manage example status page of the organization with ID 1.
resource "statuspal_status_page" "example" {
  organization_id = "1"
  status_page = {
    name      = "Example Terraform Status Page"
    url       = "example.com"
    time_zone = "Europe/Berlin"
  }
}
