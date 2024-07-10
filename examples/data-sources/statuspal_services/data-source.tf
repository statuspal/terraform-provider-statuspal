# List all services of the status page with subdomain "example-com".
data "statuspal_services" "all" {
  status_page_subdomain = "example-com"
}
