# List all status pages of the organization with ID 1.
data "statuspal_status_pages" "all" {
  organization_id = "1"
}
