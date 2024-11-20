terraform {
  required_providers {
    statuspal = {
      source = "registry.terraform.io/statuspal/statuspal"
    }
  }

  required_version = ">= 1.2.0"
}

provider "statuspal" {
  api_key = "uk_aERPQU1kUzUrRmplaXJRMlc2TDEwZz09"
  region  = "US" // "US" or "EU"
}

resource "statuspal_service" "edu" {
  status_page_subdomain = "example-com"
  service = {
    name = "Service Created from Terraform"
  }
}

output "edu_service" {
  value = statuspal_service.edu
}

resource "statuspal_service" "edu_child" {
  status_page_subdomain = statuspal_service.edu.status_page_subdomain
  service = {
    name      = "Child Service Created from Terraform"
    parent_id = statuspal_service.edu.service.id
  }
}

output "edu_child_service" {
  value = statuspal_service.edu_child
}
