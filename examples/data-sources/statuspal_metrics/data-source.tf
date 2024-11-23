# List all metrics of the status page with subdomain "example-com".
data "statuspal_metrics" "example" {
  status_page_subdomain = "example-com"
}
