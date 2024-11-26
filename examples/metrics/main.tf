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

data "statuspal_metrics" "edu" {
  status_page_subdomain = "example-com"
}


output "edu_statuspal_metrics" {
  value = data.statuspal_metrics.edu
}
