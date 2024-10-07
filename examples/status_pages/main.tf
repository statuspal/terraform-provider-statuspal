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
}

data "statuspal_status_pages" "edu" {
  organization_id = "1"
}


output "edu_status_pages" {
  value = data.statuspal_status_pages.edu
}
