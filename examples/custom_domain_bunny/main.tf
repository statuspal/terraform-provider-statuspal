terraform {
  required_providers {
    statuspal = {
      source = "registry.terraform.io/statuspal/statuspal"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 4.0"
    }
  }

  required_version = ">= 1.4.0"
}

provider "statuspal" {
  api_key = var.statuspal_api_key
  region  = "EU"
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

variable "cloudflare_zone_id" {}

variable "org_id" {}

# Step 1 — Create the status page with domain_config using Bunny CDN.
# The provider polls until the Bunny pull zone is ready and the CNAME value
# is available, so validation_records["cname"] is populated after this step.
resource "statuspal_status_page" "main" {
  organization_id = var.org_id
  status_page = {
    name      = "Acme Status"
    url       = "https://status.acme.com"
    time_zone = "UTC"
    subdomain = "acme"

    domain_config = {
      provider = "bunny"
      domain   = "status.acme.com"
    }
  }
}

# Step 2 — CNAME record to route the custom domain to the Bunny CDN hostname.
resource "cloudflare_record" "cname" {
  zone_id = var.cloudflare_zone_id
  name    = statuspal_status_page.main.status_page.domain_config.validation_records["cname"].name
  type    = "CNAME"
  content = statuspal_status_page.main.status_page.domain_config.validation_records["cname"].value
  proxied = false
  ttl     = 120
}

# Step 3 — Wait until the custom domain is fully active.
# Bunny handles SSL automatically — no statuspal_domain_ssl_records or TXT
# record is needed (unlike the Cloudflare flow).
resource "statuspal_custom_domain_validation" "main" {
  organization_id       = var.org_id
  status_page_subdomain = statuspal_status_page.main.status_page.subdomain
  timeout_seconds       = 600

  depends_on = [cloudflare_record.cname]
}

output "domain_status" {
  value = statuspal_status_page.main.status_page.domain_config.status
}
