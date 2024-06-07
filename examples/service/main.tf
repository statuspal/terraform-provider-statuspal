terraform {
  required_providers {
    statuspal = {
      source = "registry.terraform.io/statuspal/statuspal"
    }
  }
  required_version = ">= 1.1.0"
}

provider "statuspal" {
  api_key = "uk_aERPQU1kUzUrRmplaXJRMlc2TDEwZz09"
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
