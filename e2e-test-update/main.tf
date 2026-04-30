terraform {
  required_providers {
    statuspal = {
      source = "statuspal/statuspal"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 4.0"
    }
  }
}

provider "statuspal" {
  api_key = var.statuspal_api_key
  region  = "US"
}

provider "cloudflare" {
  api_token = var.cloudflare_api_token
}

variable "statuspal_api_key" { sensitive = true }
variable "cloudflare_api_token" { sensitive = true }
variable "cloudflare_zone_id" {}
variable "org_id" {}
variable "status_page_name" { default = "TF Custom Domain Test B" }
variable "status_page_url" {}
variable "subdomain" {}
variable "custom_domain" {}

# ---------------------------------------------------------------------------
# Phase 1: bare status page, no custom domain.
#   $ terraform apply
# ---------------------------------------------------------------------------
resource "statuspal_status_page" "test" {
  organization_id = var.org_id
  status_page = {
    name      = var.status_page_name
    url       = var.status_page_url
    time_zone = "UTC"
    subdomain = var.subdomain

    # ---- Phase 2: uncomment the block below and re-apply -------------------
    # domain_config = {
    #   provider = "cloudflare"
    #   domain   = var.custom_domain
    # }
    # ------------------------------------------------------------------------
  }
}

# ---------------------------------------------------------------------------
# Phase 2: uncomment everything below and re-apply.
#   $ terraform apply
# ---------------------------------------------------------------------------

# resource "cloudflare_record" "cname" {
#   zone_id = var.cloudflare_zone_id
#   name    = statuspal_status_page.test.status_page.domain_config.validation_records["cname"].name
#   type    = "CNAME"
#   content = statuspal_status_page.test.status_page.domain_config.validation_records["cname"].value
#   proxied = false
#   ttl     = 120
# }

# resource "statuspal_domain_ssl_records" "test" {
#   organization_id       = var.org_id
#   status_page_subdomain = statuspal_status_page.test.status_page.subdomain
#   timeout_seconds       = 600
#
#   depends_on = [cloudflare_record.cname]
# }

# resource "cloudflare_record" "txt" {
#   zone_id = var.cloudflare_zone_id
#   name    = statuspal_domain_ssl_records.test.certificate_txt_name
#   type    = "TXT"
#   content = statuspal_domain_ssl_records.test.certificate_txt_value
#   ttl     = 120
# }

# resource "statuspal_custom_domain_validation" "test" {
#   organization_id       = var.org_id
#   status_page_subdomain = statuspal_status_page.test.status_page.subdomain
#   timeout_seconds       = 600
#
#   depends_on = [cloudflare_record.txt]
# }

output "domain_config" {
  value = statuspal_status_page.test.status_page.domain_config
}
