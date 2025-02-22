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
  region  = "US"
}

resource "statuspal_metric" "edu" {
  status_page_subdomain = "example-com"
  metric = {
    title = "example"
    unit  = "ms"
    type  = "rt"
  }
}

output "example_statuspal_metric" {
  value = statuspal_metric.edu
}
