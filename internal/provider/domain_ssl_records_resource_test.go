package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDomainSslRecordsResource(t *testing.T) {
	var pollCount atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		count := pollCount.Add(1)

		// First 2 polls return no certificate_txt records (CNAME not yet processed).
		// From poll 3 onwards return the certificate records.
		validationRecords := map[string]string{
			"hostname_cname_name":  "example.com",
			"hostname_cname_value": "app.statuspal.tech",
		}
		if count >= 3 {
			validationRecords["certificate_txt_name"] = "_acme-challenge.example.com"
			validationRecords["certificate_txt_value"] = "abcdef1234567890"
		}

		domainStatus := "configuring"
		body := map[string]any{
			"status_page": map[string]any{
				"subdomain": "test-subdomain",
				"domain_config": map[string]any{
					"provider":           "cloudflare",
					"domain":             "example.com",
					"status":             domainStatus,
					"validation_records": validationRecords,
					"main_hostname":      nil,
					"external_id":        nil,
					"error":              nil,
					"pullzone_id":        nil,
				},
			},
		}

		if err := json.NewEncoder(w).Encode(body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainSslRecordsConfig(server.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("statuspal_domain_ssl_records.test", "certificate_txt_name", "_acme-challenge.example.com"),
					resource.TestCheckResourceAttr("statuspal_domain_ssl_records.test", "certificate_txt_value", "abcdef1234567890"),
					resource.TestCheckResourceAttr("statuspal_domain_ssl_records.test", "organization_id", "test-org"),
					resource.TestCheckResourceAttr("statuspal_domain_ssl_records.test", "status_page_subdomain", "test-subdomain"),
				),
			},
		},
	})
}

func testAccDomainSslRecordsConfig(testURL string) string {
	return fmt.Sprintf(`
provider "statuspal" {
  test_url = %q
}

resource "statuspal_domain_ssl_records" "test" {
  organization_id       = "test-org"
  status_page_subdomain = "test-subdomain"
  timeout_seconds       = 60
}
`, testURL)
}
