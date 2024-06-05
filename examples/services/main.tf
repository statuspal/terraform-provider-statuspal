terraform {
  required_providers {
    statuspal = {
      source = "registry.terraform.io/statuspal/statuspal"
    }
  }
}

provider "statuspal" {
  api_key = "uk_aERPQU1kUzUrRmplaXJRMlc2TDEwZz09"
  region  = "dev"
}

data "statuspal_services" "edu" {
  status_page_subdomain = "example-com"
}


output "edu_statuspal_services" {
  value = data.statuspal_services.edu
}
