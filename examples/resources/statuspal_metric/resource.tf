# Manage example metric of the status page with subdomain "example-com".
resource "statuspal_metric" "example" {
  status_page_subdomain = "example-com"
  metric = {
    title = "example"
    unit  = "ms"
    type  = "rt"
  }
}
