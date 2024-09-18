# Configure terraform
terraform {
  required_providers {
    statuspal = {
      source  = "statuspal/statuspal"
      version = "0.2.10"
    }
  }
}

# Configure the StatusPal provider
provider "statuspal" {
  api_key = "uk_aERPQU1kUzUrRmplaXJRMlc2TDEwZz09" // your user or organization api key
  region  = "US"                                  // "US" or "EU"
}

# List all status pages of the organization with ID 1.
data "statuspal_status_pages" "all" {
  organization_id = "1"
}

# Manage example status page of the organization with ID 1.
resource "statuspal_status_page" "example" {
  organization_id = data.statuspal_status_pages.all.organization_id // you can use it from previously defined data source
  status_page = {
    name      = "Example Terraform Status Page"
    url       = "example.com"
    time_zone = "America/New_York"
  }
}

# List all services of the previously defined example status page.
data "statuspal_services" "all" {
  status_page_subdomain = statuspal_status_page.example.status_page.subdomain
}

# Manage example service of the previously defined example status page's subdomain.
resource "statuspal_service" "example" {
  status_page_subdomain = statuspal_status_page.example.status_page.subdomain
  service = {
    name = "Example Terraform Service"
  }
}

# Manage example child service of the previously defined example service.
resource "statuspal_service" "child" {
  status_page_subdomain = statuspal_status_page.example.status_page.subdomain
  service = {
    name      = "Example Terraform Child Service"
    parent_id = statuspal_service.example.service.id
  }
}
