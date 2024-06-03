terraform {
  required_providers {
    statuspal = {
      source = "registry.terraform.io/hashicorp/statuspal"
    }
  }
  required_version = ">= 1.8.0"
}

provider "statuspal" {
  api_key = "uk_aERPQU1kUzUrRmplaXJRMlc2TDEwZz09"
  region  = "dev"
}

output "demo_sum" {
  value = provider::statuspal::demo_sum(5, 8)
}
