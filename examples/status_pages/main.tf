terraform {
  required_providers {
    statuspal = {
      source = "registry.terraform.io/hashicorp/statuspal"
    }
  }
}

provider "statuspal" {
  api_key = "uk_aERPQU1kUzUrRmplaXJRMlc2TDEwZz09"
  region  = "dev"
}

data "statuspal_status_pages" "edu" {
  organization_id = "1"
}


output "edu_status_pages" {
  value = data.statuspal_status_pages.edu
}
