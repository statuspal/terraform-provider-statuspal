package provider

import (
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceResource(t *testing.T) {
	mux := http.NewServeMux()
	basicResponseBody := `{
		"service": {
			"id": 2,
			"parent_id": 3,
			"name": "Test Service from Terraform",
			"private": false,
			"description": null,
			"monitoring": null,
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
					"name": "Test Service from Terraform",
					"description": ""
				},
				"es": {
					"name": "web ES",
					"description": ""
				},
				"fr": {
					"name": "web FR",
					"description": ""
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
			"inbound_email_id": "d346f35e-0749-4ed7-a88b-7caa679d1959"
		}
	}`
	updatedResponseBody := strings.Replace(basicResponseBody, `"name": "Test Service from Terraform"`, `"name": "Edited Test Service from Terraform"`, 2)
	updatedResponseBody = strings.Replace(updatedResponseBody, `"updated_at": "2024-05-16T10:00:00"`, `"updated_at": "2024-05-20T10:00:00"`, 1)

	// Mock create response for resource
	mux.HandleFunc("/status_pages/terraform-test/services", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(basicResponseBody)); err != nil {
			log.Printf(`Error writing "/status_pages/terraform-test/services" response with method "%s": %v`, r.Method, err)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})

	// Mock after create read response for resource
	mux.HandleFunc("/status_pages/terraform-test/services/2", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(basicResponseBody)); err != nil {
			log.Printf(`Error writing "/status_pages/terraform-test/services/2" response with method "%s": %v`, r.Method, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	// Mock update, read and delete responses for resource
	mux.HandleFunc("/status_pages/terraform-test-updated/services/2", func(w http.ResponseWriter, r *http.Request) {
		responseBody := basicResponseBody

		switch r.Method {
		case http.MethodGet:
			responseBody = updatedResponseBody
		case http.MethodPut:
			responseBody = updatedResponseBody
		case http.MethodDelete:
			responseBody = `""`
		}

		if _, err := w.Write([]byte(responseBody)); err != nil {
			log.Printf(`Error writing "/status_pages/terraform-test-updated/services/2" response with method "%s": %v`, r.Method, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	childResponseBody := strings.Replace(basicResponseBody, `"id": 2`, `"id": 5`, 1)
	childResponseBody = strings.Replace(childResponseBody, `"name": "Test Service from Terraform"`, `"name": "Test Child Service from Terraform"`, 2)
	childResponseBody = strings.Replace(childResponseBody, `"parent_id": 3`, `"parent_id": 2`, 1)

	// Mock create response for parent resource
	mux.HandleFunc("/status_pages/terraform-test-parent/services", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(basicResponseBody)); err != nil {
			log.Printf(`Error writing "/status_pages/terraform-test-parent/services" response with method "%s": %v`, r.Method, err)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})

	// Mock read and delete responses for parent resource
	mux.HandleFunc("/status_pages/terraform-test-parent/services/2", func(w http.ResponseWriter, r *http.Request) {
		responseBody := basicResponseBody

		switch r.Method {
		case http.MethodDelete:
			responseBody = `""`
		}

		if _, err := w.Write([]byte(responseBody)); err != nil {
			log.Printf(`Error writing "/status_pages/terraform-test-parent/services/2" response with method "%s": %v`, r.Method, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	// Mock create response for child resource
	mux.HandleFunc("/status_pages/terraform-test-child/services", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(childResponseBody)); err != nil {
			log.Printf(`Error writing "/status_pages/terraform-test-child/services" response with method "%s": %v`, r.Method, err)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})

	// Mock read and delete responses for child resource
	mux.HandleFunc("/status_pages/terraform-test-child/services/5", func(w http.ResponseWriter, r *http.Request) {
		responseBody := childResponseBody

		switch r.Method {
		case http.MethodDelete:
			responseBody = `""`
		}

		if _, err := w.Write([]byte(responseBody)); err != nil {
			log.Printf(`Error writing "/status_pages/terraform-test-child/services/5" response with method "%s": %v`, r.Method, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mockServer := httptest.NewServer(mux)
	defer mockServer.Close()
	providerConfig := providerConfig(&mockServer.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Missing required status_page_subdomain attribute error testing
			{
				Config: *providerConfig + `resource "statuspal_service" "test" {
					service = {}
				}`,
				ExpectError: regexp.MustCompile(`The argument "status_page_subdomain" is required, but no definition was\nfound.`),
			},
			// Missing a required service attribute error testing
			{
				Config: *providerConfig + `resource "statuspal_service" "test" {
					status_page_subdomain = "terraform-test"
					service = {}
				}`,
				ExpectError: regexp.MustCompile(`Inappropriate value for attribute "service": attribute "name" is required.`),
			},
			// Create and Read testing
			{
				Config: *providerConfig + `resource "statuspal_service" "test" {
					status_page_subdomain = "terraform-test"
					service = {
						name = "Test Service from Terraform"
						translations = {
							en = {
								name = "Test Service from Terraform"
								description = ""
							}
							es = {
								name = "web ES"
								description = ""
							}
							fr = {
								name = "web FR"
								description = ""
							}
						}
						webhook_custom_jsonpath_settings = {
							jsonpath = ""
							expected_result = ""
						}
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspal_service.test", "status_page_subdomain", "terraform-test"),
					// Verify service
					resource.TestCheckResourceAttr("statuspal_service.test", "service.%", "27"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.id", "2"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.name", "Test Service from Terraform"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.description", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.private_description", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.parent_id", "3"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.current_incident_type", "custom-type"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.monitoring", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.webhook_monitoring_service", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.webhook_custom_jsonpath_settings.%", "2"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.webhook_custom_jsonpath_settings.jsonpath", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.webhook_custom_jsonpath_settings.expected_result", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.inbound_email_address", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.incoming_webhook_url", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.ping_url", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.incident_type", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.parent_incident_type", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.is_up", "false"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.pause_monitoring_during_maintenances", "false"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.inbound_email_id", "d346f35e-0749-4ed7-a88b-7caa679d1959"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.auto_incident", "false"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.auto_notify", "false"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.children_ids.#", "2"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.children_ids.0", "343"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.children_ids.1", "656"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.translations.%", "3"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.translations.en.%", "2"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.translations.en.name", "Test Service from Terraform"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.translations.en.description", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.translations.es.%", "2"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.translations.es.name", "web ES"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.translations.es.description", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.translations.fr.%", "2"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.translations.fr.name", "web FR"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.translations.fr.description", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.private", "false"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.display_uptime_graph", "false"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.display_response_time_chart", "false"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.order", "3"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.inserted_at", "2023-11-15T10:03:20"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.updated_at", "2024-05-16T10:00:00"),
					// Verify placeholder id attribute
					resource.TestCheckResourceAttr("statuspal_service.test", "id", "placeholder"),
				),
			},
			// ImportState fail testing
			{
				ResourceName:      "statuspal_service.test",
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile(`Expected StatusPal service import identifier with format:\n"<status_page_subdomain> <service_id>"`),
			},
			// ImportState testing
			{
				ResourceName:      "statuspal_service.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "terraform-test 2",
			},
			// Update and Read testing
			{
				Config: *providerConfig + `resource "statuspal_service" "test" {
					status_page_subdomain = "terraform-test-updated"
					service = {
						name = "Edited Test Service from Terraform"
						translations = {
							en = {
								name = "Edited Test Service from Terraform"
								description = ""
							}
							es = {
								name = "web ES"
								description = ""
							}
							fr = {
								name = "web FR"
								description = ""
							}
						}
						webhook_custom_jsonpath_settings = {
							jsonpath = ""
							expected_result = ""
						}
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspal_service.test", "status_page_subdomain", "terraform-test-updated"),
					// Verify service
					resource.TestCheckResourceAttr("statuspal_service.test", "service.%", "27"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.id", "2"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.name", "Edited Test Service from Terraform"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.description", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.private_description", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.parent_id", "3"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.current_incident_type", "custom-type"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.monitoring", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.webhook_monitoring_service", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.webhook_custom_jsonpath_settings.%", "2"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.webhook_custom_jsonpath_settings.jsonpath", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.webhook_custom_jsonpath_settings.expected_result", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.inbound_email_address", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.incoming_webhook_url", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.ping_url", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.incident_type", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.parent_incident_type", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.is_up", "false"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.pause_monitoring_during_maintenances", "false"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.inbound_email_id", "d346f35e-0749-4ed7-a88b-7caa679d1959"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.auto_incident", "false"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.auto_notify", "false"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.children_ids.#", "2"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.children_ids.0", "343"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.children_ids.1", "656"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.translations.%", "3"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.translations.en.%", "2"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.translations.en.name", "Edited Test Service from Terraform"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.translations.en.description", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.translations.es.%", "2"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.translations.es.name", "web ES"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.translations.es.description", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.translations.fr.%", "2"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.translations.fr.name", "web FR"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.translations.fr.description", ""),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.private", "false"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.display_uptime_graph", "false"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.display_response_time_chart", "false"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.order", "3"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.inserted_at", "2023-11-15T10:03:20"),
					resource.TestCheckResourceAttr("statuspal_service.test", "service.updated_at", "2024-05-20T10:00:00"),
					// Verify placeholder id attribute
					resource.TestCheckResourceAttr("statuspal_service.test", "id", "placeholder"),
				),
			},
			// Creating a child service
			{
				Config: *providerConfig + `resource "statuspal_service" "parent_test" {
					status_page_subdomain = "terraform-test-parent"
					service = {
						name = "Test Service from Terraform"
						translations = {
							en = {
								name = "Test Service from Terraform"
								description = ""
							}
							es = {
								name = "web ES"
								description = ""
							}
							fr = {
								name = "web FR"
								description = ""
							}
						}
						webhook_custom_jsonpath_settings = {
							jsonpath = ""
							expected_result = ""
						}
					}
				}

				resource "statuspal_service" "child_test" {
					status_page_subdomain = "terraform-test-child"
					service = {
						name = "Test Child Service from Terraform"
						parent_id = statuspal_service.parent_test.service.id
						translations = {
							en = {
								name = "Test Child Service from Terraform"
								description = ""
							}
							es = {
								name = "web ES"
								description = ""
							}
							fr = {
								name = "web FR"
								description = ""
							}
						}
						webhook_custom_jsonpath_settings = {
							jsonpath = ""
							expected_result = ""
						}
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspal_service.child_test", "service.parent_id", "2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
