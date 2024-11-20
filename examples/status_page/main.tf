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

resource "statuspal_status_page" "edu" {
  organization_id = "1"
  status_page = {
    name      = "Status Page Created from Terraform"
    url       = "terraform-test.com"
    time_zone = "Europe/Berlin"
  }
}

output "edu_status_page" {
  value = statuspal_status_page.edu
}
