package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDomainSslRecordsResource(t *testing.T) {
	var pollCount atomic.Int32
	var activated atomic.Bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// During Create:
		//   - First 2 polls return no certificate_txt records (CNAME not yet processed).
		//   - From poll 3 onwards return the certificate records.
		// After Create completes (activated flag set by the test), simulate
		// post-activation API: status flips to "active" and the API stops
		// returning the certificate_txt_* keys. Read must preserve the
		// existing state values rather than clearing them.
		validationRecords := map[string]string{
			"hostname_cname_name":  "example.com",
			"hostname_cname_value": "app.statuspal.tech",
		}
		domainStatus := "configuring"

		if activated.Load() {
			domainStatus = "active"
			// Post-activation: API still includes the keys but with empty
			// string values. Read must treat this the same as "absent" and
			// preserve the existing state values.
			validationRecords["certificate_txt_name"] = ""
			validationRecords["certificate_txt_value"] = ""
		} else {
			count := pollCount.Add(1)
			if count >= 3 {
				validationRecords["certificate_txt_name"] = "_acme-challenge.example.com"
				validationRecords["certificate_txt_value"] = "abcdef1234567890"
			}
		}

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
					func(_ *terraform.State) error {
						activated.Store(true)
						return nil
					},
				),
			},
			// Refresh-only step after the API has flipped to post-activation mode.
			// The certificate_txt_* values must survive Read intact, otherwise the
			// downstream cloudflare_record.txt would be forced to replace.
			{
				RefreshState:       true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("statuspal_domain_ssl_records.test", "certificate_txt_name", "_acme-challenge.example.com"),
					resource.TestCheckResourceAttr("statuspal_domain_ssl_records.test", "certificate_txt_value", "abcdef1234567890"),
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
