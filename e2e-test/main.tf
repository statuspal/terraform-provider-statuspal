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

variable "statuspal_api_key" {
  sensitive = true
}

variable "cloudflare_api_token" {
  sensitive = true
}

variable "cloudflare_zone_id" {
  description = "Cloudflare zone ID for the domain you're testing against."
}

variable "org_id" {
  description = "StatusPal organization ID."
}

variable "status_page_name" {
  description = "Display name of the status page."
  default     = "Terraform E2E Test"
}

variable "status_page_url" {
  description = "Public URL of the company/site this status page is for."
}

variable "subdomain" {
  description = "Status page subdomain on StatusPal."
}

variable "custom_domain" {
  description = "Custom hostname to point at the status page (must be lowercase)."
}

resource "statuspal_status_page" "test" {
  organization_id = var.org_id
  status_page = {
    name      = var.status_page_name
    url       = var.status_page_url
    time_zone = "UTC"
    subdomain = var.subdomain

    domain_config = {
      provider = "cloudflare"
      domain   = var.custom_domain
    }
  }
}

# Step 1 — CNAME record (values available immediately from status page)
resource "cloudflare_record" "cname" {
  zone_id = var.cloudflare_zone_id
  name    = statuspal_status_page.test.status_page.domain_config.validation_records["cname"].name
  type    = "CNAME"
  content = statuspal_status_page.test.status_page.domain_config.validation_records["cname"].value
  proxied = false
  ttl     = 120
}

# Step 2 — polls until Cloudflare generates the ACME TXT challenge (requires CNAME in DNS first)
resource "statuspal_domain_ssl_records" "test" {
  organization_id       = var.org_id
  status_page_subdomain = statuspal_status_page.test.status_page.subdomain
  timeout_seconds       = 600

  depends_on = [cloudflare_record.cname]
}

# Step 3 — TXT record for SSL certificate (values computed after Step 2)
resource "cloudflare_record" "txt" {
  zone_id = var.cloudflare_zone_id
  name    = statuspal_domain_ssl_records.test.certificate_txt_name
  type    = "TXT"
  content = statuspal_domain_ssl_records.test.certificate_txt_value
  ttl     = 120
}

# Step 4 — waiter, blocks until domain is active
resource "statuspal_custom_domain_validation" "test" {
  organization_id       = var.org_id
  status_page_subdomain = statuspal_status_page.test.status_page.subdomain
  timeout_seconds       = 600

  depends_on = [cloudflare_record.txt]
}

output "domain_config" {
  value = statuspal_status_page.test.status_page.domain_config
}

output "validation_records" {
  value = statuspal_status_page.test.status_page.domain_config.validation_records
}
