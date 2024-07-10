# Manage example service of the status page with subdomain "example-com".
resource "statuspal_service" "example" {
  status_page_subdomain = "example-com"
  service = {
    name = "Example Terraform Service"
  }
}

# Manage example child service of the previously defined example service.
resource "statuspal_service" "child" {
  status_page_subdomain = statuspal_service.example.status_page_subdomain
  service = {
    name      = "Example Terraform Child Service"
    parent_id = statuspal_service.example.service.id
  }
}
