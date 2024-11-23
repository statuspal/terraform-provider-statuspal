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

resource "statuspal_status_page" "example" {
  organization_id = "1"
  status_page = {
    name      = "example"
    url       = "example.com"
    time_zone = "Europe/Berlin"
  }
}

resource "statuspal_metric" "example" {
  status_page_subdomain = statuspal_status_page.example.status_page.subdomain
  metric = {
    title = "example"
    unit  = "ms"
    type  = "rt"
  }
}

data "statuspal_metrics" "example" {
  status_page_subdomain = statuspal_status_page.example.status_page.subdomain
}

output "example_statuspal_metrics" {
  value = data.statuspal_metrics.example.metrics[0].title
}
