package provider

import (
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServicesDataSource(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/status_pages/terraform-test/services", func(w http.ResponseWriter, r *http.Request) {
		// Mock response for data source
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{
			"links": {
				"next": null,
				"prev": null
			},
			"meta": {
				"total_count": 3
			},
			"services": [
				{
					"id": 2,
					"parent_id": 3,
					"name": "api",
					"private": false,
					"description": null,
					"monitoring": "3rd_party",
					"webhook_monitoring_service": null,
					"webhook_custom_jsonpath_settings": null,
					"inbound_email_address": null,
					"incoming_webhook_url": null,
					"inserted_at": "2023-11-15T10:03:20",
					"updated_at": "2024-05-16T10:00:00",
					"order": 3,
					"incident_type": null,
					"translations": {
						"en": {
							"description": "",
							"name": "web EN"
						},
						"es": {
							"description": "",
							"name": "web ES"
						},
						"fr": {
							"description": "",
							"name": "web FR"
						}
					},
					"auto_notify": false,
					"current_incident_type": "custom-type",
					"parent_incident_type": null,
					"children_ids": [
						343,
						656
					],
					"is_up": null,
					"auto_incident": false,
					"ping_url": null,
					"pause_monitoring_during_maintenances": false,
					"private_description": null,
					"display_response_time_chart": false,
					"display_uptime_graph": false,
					"inbound_email_id": "d346f35e-0749-4ed7-a88b-7caa679d1959",
					"monitoring_options": {
						"headers": [],
						"method": "",
						"keyword_down": "DOWN",
						"keyword_up": "UP"
					}					
				},
				{
					"id": 1,
					"parent_id": null,
					"name": "web FR",
					"private": false,
					"description": "",
					"monitoring": "webhook",
					"webhook_monitoring_service": "custom-jsonpath",
					"webhook_custom_jsonpath_settings": {
						"jsonpath": "$.status",
						"expected_result": "\"up\""
					},
					"inbound_email_address": null,
					"incoming_webhook_url": null,
					"inserted_at": "2023-11-15T10:03:20",
					"updated_at": "2024-04-09T11:20:05",
					"order": 2,
					"incident_type": null,
					"translations": {
						"en": {
							"description": "",
							"name": "web EN"
						},
						"es": {
							"description": "",
							"name": "web EN"
						},
						"fr": {
							"description": "",
							"name": "web FR"
						}
					},
					"auto_notify": true,
					"current_incident_type": null,
					"parent_incident_type": null,
					"children_ids": [],
					"is_up": false,
					"auto_incident": true,
					"ping_url": "https://local.statuspal.io:4001/api/v2/status_pages/pontsystems-eu-hu/services/14be3b68-1f77-4732-a49a-05eeea5515de/automate/custom-jsonpath",
					"pause_monitoring_during_maintenances": false,
					"private_description": null,
					"display_response_time_chart": false,
					"display_uptime_graph": false,
					"inbound_email_id": "14be3b68-1f77-4732-a49a-05eeea5515de"
				},
				{
					"id": 4,
					"parent_id": null,
					"name": "new service",
					"private": true,
					"description": "",
					"monitoring": "internal",
					"webhook_monitoring_service": null,
					"webhook_custom_jsonpath_settings": null,
					"inbound_email_address": null,
					"incoming_webhook_url": null,
					"inserted_at": "2023-12-13T09:42:54",
					"updated_at": "2024-04-09T11:20:05",
					"order": 3,
					"incident_type": null,
					"translations": {
						"en": {
							"description": "",
							"name": "new service"
						},
						"fr": {
							"description": "",
							"name": ""
						}
					},
					"auto_notify": false,
					"current_incident_type": null,
					"parent_incident_type": null,
					"children_ids": [],
					"is_up": null,
					"auto_incident": false,
					"ping_url": null,
					"pause_monitoring_during_maintenances": false,
					"private_description": null,
					"display_response_time_chart": false,
					"display_uptime_graph": false,
					"inbound_email_id": "ca237edd-72ac-4793-b278-5e682e3d7b47",
					"monitoring_options": {
						"headers": [
							{
								"value": "abcdef",
								"key": "Authorization"
							},
							{
								"value": "es",
								"key": "accept-language"
							}
						],
						"method": "head",
						"keyword_down": null,
						"keyword_up": null
					}
				}
			]
		}`)); err != nil {
			log.Printf(`Error writing "/status_pages/terraform-test/services" response: %v`, err)
			return
		}
	})
	mockServer := httptest.NewServer(mux)
	defer mockServer.Close()
	providerConfig := providerConfig(&mockServer.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Missing required status_page_subdomain error testing
			{
				Config:      *providerConfig + `data "statuspal_services" "test" {}`,
				ExpectError: regexp.MustCompile(`The argument "status_page_subdomain" is required, but no definition was\nfound.`),
			},
			// Read testing
			{
				Config: *providerConfig + `data "statuspal_services" "test" {
					status_page_subdomain = "terraform-test"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.statuspal_services.test", "status_page_subdomain", "terraform-test"),
					// Verify number of services returned
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.#", "3"),
					// Verify the first service to ensure all attributes are set
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.%", "28"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.id", "2"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.name", "api"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.description", ""),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.private_description", ""),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.parent_id", "3"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.current_incident_type", "custom-type"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.monitoring", "3rd_party"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.webhook_monitoring_service", ""),
					resource.TestCheckNoResourceAttr("data.statuspal_services.test", "services.0.webhook_custom_jsonpath_settings"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.inbound_email_address", ""),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.incoming_webhook_url", ""),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.ping_url", ""),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.incident_type", ""),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.parent_incident_type", ""),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.is_up", "false"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.pause_monitoring_during_maintenances", "false"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.inbound_email_id", "d346f35e-0749-4ed7-a88b-7caa679d1959"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.auto_incident", "false"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.auto_notify", "false"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.children_ids.#", "2"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.children_ids.0", "343"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.children_ids.1", "656"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.translations.%", "3"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.translations.en.%", "2"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.translations.en.name", "web EN"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.translations.en.description", ""),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.translations.es.%", "2"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.translations.es.name", "web ES"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.translations.es.description", ""),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.translations.fr.%", "2"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.translations.fr.name", "web FR"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.translations.fr.description", ""),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.private", "false"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.display_uptime_graph", "false"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.display_response_time_chart", "false"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.order", "3"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.inserted_at", "2023-11-15T10:03:20"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.updated_at", "2024-05-16T10:00:00"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.monitoring_options.keyword_up", "UP"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.0.monitoring_options.keyword_down", "DOWN"),
					// Verify the second service webhook_custom_jsonpath_settings attribute
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.1.webhook_custom_jsonpath_settings.%", "2"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.1.webhook_custom_jsonpath_settings.jsonpath", "$.status"),
					resource.TestCheckResourceAttr("data.statuspal_services.test", "services.1.webhook_custom_jsonpath_settings.expected_result", "\"up\""),
					// Verify placeholder id attribute
					resource.TestCheckResourceAttr("data.statuspal_services.test", "id", "placeholder"),
				),
			},
		},
	})
}
