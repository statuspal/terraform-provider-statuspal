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

resource "statuspal_status_page" "test" {
  organization_id = "5883"
  status_page = {
    name      = "Claude Terraform E2E Test"
    url       = "https://messuti.io"
    time_zone = "UTC"
    subdomain = "claude-tf-e2e-4"

    domain_config = {
      provider = "cloudflare"
      domain   = "claude-terraform-test-4.messuti.io"
    }
  }
}

# Step 1 — CNAME record (values available immediately from status page)
resource "cloudflare_record" "cname" {
  zone_id = "aba95479eb6a6e70b7e1224ed4c7f2e3"
  name    = statuspal_status_page.test.status_page.domain_config.validation_records["cname"].name
  type    = "CNAME"
  content = statuspal_status_page.test.status_page.domain_config.validation_records["cname"].value
  proxied = false
  ttl     = 120
}

# Step 2 — polls until Cloudflare generates the ACME TXT challenge (requires CNAME in DNS first)
resource "statuspal_domain_ssl_records" "test" {
  organization_id       = "5883"
  status_page_subdomain = statuspal_status_page.test.status_page.subdomain
  timeout_seconds       = 600

  depends_on = [cloudflare_record.cname]
}

# Step 3 — TXT record for SSL certificate (values computed after Step 2)
resource "cloudflare_record" "txt" {
  zone_id = "aba95479eb6a6e70b7e1224ed4c7f2e3"
  name    = statuspal_domain_ssl_records.test.certificate_txt_name
  type    = "TXT"
  content = statuspal_domain_ssl_records.test.certificate_txt_value
  ttl     = 120
}

# Step 4 — waiter, blocks until domain is active
resource "statuspal_custom_domain_validation" "test" {
  organization_id       = "5883"
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
