terraform {
  required_providers {
    statuspal = {
      source = "registry.terraform.io/statuspal/statuspal"
    }
  }
}

provider "statuspal" {}

data "statuspal_status_pages" "example" {}
