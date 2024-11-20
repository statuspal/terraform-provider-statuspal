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

data "statuspal_services" "edu" {
  status_page_subdomain = "example-com"
}


output "edu_statuspal_services" {
  value = data.statuspal_services.edu
}
