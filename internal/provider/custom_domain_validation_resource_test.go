package provider

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomDomainValidationResource(t *testing.T) {
	// pollCount tracks how many GET calls the waiter has made so we can
	// simulate a domain that starts as "configuring" and later becomes "active".
	var pollCount atomic.Int32

	statusPageBody := func(status string) string {
		return fmt.Sprintf(`{
			"status_page": {
				"name": "Test Status Page",
				"subdomain": "terraform-test",
				"url": "terraform.test",
				"time_zone": "UTC",
				"domain_config": {
					"provider": "cloudflare",
					"domain": "status.terraform.test",
					"main_hostname": "ssl-for-saas.example.com",
					"status": %q,
					"error": null,
					"external_id": "ext-abc123",
					"pullzone_id": null,
					"validation_records": {
						"hostname_cname_name": "status.terraform.test",
						"hostname_cname_value": "ssl-for-saas.example.com",
						"hostname_txt_name": "_cf-custom-hostname.status.terraform.test",
						"hostname_txt_value": "some-verification-token"
					}
				}
			}
		}`, status)
	}

	mux := http.NewServeMux()

	// Status page endpoint: first 2 polls return "configuring", then "active".
	mux.HandleFunc("/orgs/1/status_pages/terraform-test", func(w http.ResponseWriter, r *http.Request) {
		var body string
		count := pollCount.Add(1)
		if count <= 2 {
			body = statusPageBody("configuring")
		} else {
			body = statusPageBody("active")
		}
		if _, err := w.Write([]byte(body)); err != nil {
			log.Printf("Error writing status page response: %v", err)
		}
	})

	mockServer := httptest.NewServer(mux)
	defer mockServer.Close()
	providerCfg := providerConfig(&mockServer.URL)

	config := *providerCfg + `
		resource "statuspal_custom_domain_validation" "test" {
			organization_id       = "1"
			status_page_subdomain = "terraform-test"
			timeout_seconds       = 60
		}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Missing required attribute error
			{
				Config:      *providerCfg + `resource "statuspal_custom_domain_validation" "test" {}`,
				ExpectError: regexp.MustCompile(`The argument "organization_id" is required`),
			},
			// Create: waiter polls through "configuring" and resolves to "active"
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspal_custom_domain_validation.test", "organization_id", "1"),
					resource.TestCheckResourceAttr("statuspal_custom_domain_validation.test", "status_page_subdomain", "terraform-test"),
					resource.TestCheckResourceAttr("statuspal_custom_domain_validation.test", "timeout_seconds", "60"),
					resource.TestCheckResourceAttr("statuspal_custom_domain_validation.test", "id", "placeholder"),
				),
			},
		},
	})
}

func TestAccCustomDomainValidationResource_FailedToConfigure(t *testing.T) {
	failedBody := `{
		"status_page": {
			"name": "Test Status Page",
			"subdomain": "terraform-test-fail",
			"url": "terraform.test",
			"time_zone": "UTC",
			"domain_config": {
				"provider": "cloudflare",
				"domain": "status.terraform.test",
				"main_hostname": null,
				"status": "failed_to_configure",
				"error": "DNS verification failed: CNAME record not found",
				"external_id": null,
				"pullzone_id": null,
				"validation_records": null
			}
		}
	}`

	mux := http.NewServeMux()
	mux.HandleFunc("/orgs/1/status_pages/terraform-test-fail", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(failedBody)); err != nil {
			log.Printf("Error writing status page response: %v", err)
		}
	})

	mockServer := httptest.NewServer(mux)
	defer mockServer.Close()
	providerCfg := providerConfig(&mockServer.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: *providerCfg + `
					resource "statuspal_custom_domain_validation" "test" {
						organization_id       = "1"
						status_page_subdomain = "terraform-test-fail"
						timeout_seconds       = 60
					}
				`,
				ExpectError: regexp.MustCompile(`CNAME record not found`),
			},
		},
	})
}
