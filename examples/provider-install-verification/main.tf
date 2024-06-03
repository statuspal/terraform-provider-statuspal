terraform {
  required_providers {
    statuspal = {
      source = "registry.terraform.io/hashicorp/statuspal"
    }
  }
}

provider "statuspal" {}

data "statuspal_status_pages" "example" {}
