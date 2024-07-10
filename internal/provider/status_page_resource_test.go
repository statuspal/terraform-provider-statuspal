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

func TestAccStatusPageResource(t *testing.T) {
	mux := http.NewServeMux()
	responseBody := `{
		"status_page": {
			"theme_selected": "default",
			"scheduled_maintenance_days": 7,
			"display_uptime_graph": true,
			"inserted_at": "2024-04-15T11:20:35",
			"updated_at": "2024-04-20T11:22:32",
			"header_fg_color": "ffffff",
			"history_limit_days": 90,
			"head_code": null,
			"domain": null,
			"support_email": null,
			"locked_when_maintenance": false,
			"organization_id": 1,
			"custom_footer": null,
			"custom_incident_types_enabled": false,
			"slack_subscriptions_enabled": false,
			"date_format": null,
			"maintenance_notification_hours": 6,
			"subdomain": "terraform-test",
			"twitter_public_screen_name": null,
			"header_logo_text": "Test Status Page from Terraform",
			"member_restricted": false,
			"url": "terraform.test",
			"status_ok_color": "48CBA5",
			"uptime_graph_days": 90,
			"subscribers_enabled": true,
			"display_about": false,
			"translations": {
				"en": {
					"header_logo_text": "",
					"public_company_name": "Public company name EN"
				},
				"fr": {
					"header_logo_text": "",
					"public_company_name": "Public company name FR"
				}
			},
			"tweet_by_default": false,
			"display_calendar": true,
			"email_templates_enabled": false,
			"google_calendar_enabled": false,
			"link_color": "0c91c3",
			"email_layout_template": null,
			"name": "Test Status Page from Terraform",
			"status_major_color": "e75a53",
			"custom_header": null,
			"date_format_enforce_everywhere": false,
			"time_format": null,
			"header_bg_color1": "009688",
			"incident_link_color": null,
			"bg_image": null,
			"logo": null,
			"favicon": null,
			"custom_css": null,
			"current_incidents_position": "below_services",
			"custom_js": null,
			"minor_notification_hours": 6,
			"mattermost_notifications_enabled": false,
			"info_notices_enabled": true,
			"captcha_enabled": true,
			"about": null,
			"google_chat_notifications_enabled": false,
			"discord_notifications_enabled": false,
			"theme_configs": {
				"header_bg_color1": "",
				"header_bg_color2": "",
				"header_fg_color": "",
				"incident_header_color": "",
				"link_color": "",
				"status_maintenance_color": "",
				"status_major_color": "",
				"status_minor_color": "",
				"status_ok_color": ""
			},
			"status_minor_color": "FFA500",
			"tweeting_enabled": true,
			"sms_notifications_enabled": false,
			"notify_by_default": false,
			"hide_watermark": false,
			"custom_domain_enabled": false,
			"enable_auto_translations": false,
			"restricted_ips": null,
			"feed_enabled": true,
			"header_bg_color2": "0c91c3",
			"public_company_name": "Public company name EN",
			"notification_email": null,
			"email_notification_template": null,
			"teams_notifications_enabled": false,
			"status_maintenance_color": "5378c1",
			"email_confirmation_template": null,
			"time_zone": "Europe/Budapest",
			"calendar_enabled": false,
			"major_notification_hours": 3,
			"incident_header_color": "009688",
			"reply_to_email": null,
			"noindex": false
		}
	}`
	updatedResponseBody := strings.Replace(responseBody, `"name": "Test Status Page from Terraform"`, `"name": "Edited Test Status Page from Terraform"`, 1)
	updatedResponseBody = strings.Replace(updatedResponseBody, `"subdomain": "terraform-test"`, `"subdomain": "terraform-test-updated"`, 1)
	updatedResponseBody = strings.Replace(updatedResponseBody, `"updated_at": "2024-04-20T11:22:32"`, `"updated_at": "2024-04-25T11:22:32"`, 1)

	// Mock create response for resource
	mux.HandleFunc("/orgs/1/status_pages", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(responseBody)); err != nil {
			log.Printf(`Error writing "/orgs/1/status_pages" response with method "%s": %v`, r.Method, err)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})

	// Mock update response for resource
	mux.HandleFunc("/orgs/1/status_pages/terraform-test", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			responseBody = updatedResponseBody
		}

		if _, err := w.Write([]byte(responseBody)); err != nil {
			log.Printf(`Error writing "/orgs/1/status_pages/terraform-test" response with method "%s": %v`, r.Method, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	// Mock after update read and delete responses for resource
	mux.HandleFunc("/orgs/1/status_pages/terraform-test-updated", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			responseBody = updatedResponseBody
		case http.MethodDelete:
			responseBody = `""`
		}

		if _, err := w.Write([]byte(responseBody)); err != nil {
			log.Printf(`Error writing "/orgs/1/status_pages/terraform-test-updated" response with method "%s": %v`, r.Method, err)
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
			// Missing required organization_id attribute error testing
			{
				Config:      *providerConfig + `resource "statuspal_status_page" "test" {}`,
				ExpectError: regexp.MustCompile(`The argument "organization_id" is required, but no definition was found.`),
			},
			// Missing a required status_page attribute error testing
			{
				Config: *providerConfig + `resource "statuspal_status_page" "test" {
					organization_id = "1"
					status_page = {
						name = "Test Status Page from Terraform"
						url = "terraform.test"
					}
				}`,
				ExpectError: regexp.MustCompile(`Inappropriate value for attribute "status_page": attribute "time_zone" is\nrequired.`),
			},
			// Create and Read testing
			{
				Config: *providerConfig + `resource "statuspal_status_page" "test" {
					organization_id = "1"
					status_page = {
						name = "Test Status Page from Terraform"
						url = "terraform.test"
						time_zone = "Europe/Budapest"
						translations = {
							en = {
								header_logo_text = ""
								public_company_name = "Public company name EN"
							}
							fr = {
								header_logo_text = ""
								public_company_name = "Public company name FR"
							}
						}
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspal_status_page.test", "organization_id", "1"),
					// Verify status_page
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.%", "75"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.theme_selected", "default"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.scheduled_maintenance_days", "7"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.display_uptime_graph", "true"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.header_fg_color", "ffffff"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.history_limit_days", "90"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.head_code", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.restricted_ips", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.support_email", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.locked_when_maintenance", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.custom_footer", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.custom_incident_types_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.slack_subscriptions_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.date_format", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.maintenance_notification_hours", "6"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.subdomain", "terraform-test"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.twitter_public_screen_name", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.header_logo_text", "Test Status Page from Terraform"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.member_restricted", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.url", "terraform.test"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.status_ok_color", "48CBA5"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.uptime_graph_days", "90"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.subscribers_enabled", "true"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.display_about", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.%", "2"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.en.%", "2"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.en.header_logo_text", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.en.public_company_name", "Public company name EN"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.fr.%", "2"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.fr.header_logo_text", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.fr.public_company_name", "Public company name FR"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.tweet_by_default", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.display_calendar", "true"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.email_templates_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.google_calendar_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.link_color", "0c91c3"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.email_layout_template", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.name", "Test Status Page from Terraform"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.status_major_color", "e75a53"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.custom_header", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.date_format_enforce_everywhere", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.time_format", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.header_bg_color1", "009688"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.incident_link_color", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.bg_image", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.logo", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.favicon", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.custom_css", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.current_incidents_position", "below_services"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.custom_js", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.minor_notification_hours", "6"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.mattermost_notifications_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.info_notices_enabled", "true"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.captcha_enabled", "true"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.about", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.google_chat_notifications_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.discord_notifications_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.theme_configs.%", "9"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.theme_configs.link_color", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.theme_configs.header_bg_color1", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.theme_configs.header_bg_color2", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.theme_configs.header_fg_color", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.theme_configs.incident_header_color", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.theme_configs.status_ok_color", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.theme_configs.status_minor_color", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.theme_configs.status_major_color", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.theme_configs.status_maintenance_color", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.status_minor_color", "FFA500"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.tweeting_enabled", "true"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.sms_notifications_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.notify_by_default", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.hide_watermark", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.custom_domain_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.enable_auto_translations", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.feed_enabled", "true"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.header_bg_color2", "0c91c3"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.public_company_name", "Public company name EN"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.notification_email", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.email_notification_template", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.teams_notifications_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.status_maintenance_color", "5378c1"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.email_confirmation_template", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.time_zone", "Europe/Budapest"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.calendar_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.major_notification_hours", "3"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.incident_header_color", "009688"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.reply_to_email", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.noindex", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.inserted_at", "2024-04-15T11:20:35"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.updated_at", "2024-04-20T11:22:32"),
					// Verify placeholder id attribute
					resource.TestCheckResourceAttr("statuspal_status_page.test", "id", "placeholder"),
				),
			},
			// ImportState fail testing
			{
				ResourceName:      "statuspal_status_page.test",
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile(`Expected StatusPal status page import identifier with format:\n"<organization_id> <status_page_subdomain>"`),
			},
			// ImportState testing
			{
				ResourceName:      "statuspal_status_page.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "1 terraform-test",
				// // The last_updated attribute does not exist in the StatusPal
				// // API, therefore there is no value for it during import.
				// ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: *providerConfig + `resource "statuspal_status_page" "test" {
						organization_id = "1"
						status_page = {
							name = "Edited Test Status Page from Terraform"
							url = "terraform.test"
							time_zone = "Europe/Budapest"
							subdomain = "terraform-test-updated"
							translations = {
								en = {
									header_logo_text = ""
									public_company_name = "Public company name EN"
								}
								fr = {
									header_logo_text = ""
									public_company_name = "Public company name FR"
								}
							}
						}
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspal_status_page.test", "organization_id", "1"),
					// Verify status_page
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.%", "75"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.theme_selected", "default"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.scheduled_maintenance_days", "7"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.display_uptime_graph", "true"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.header_fg_color", "ffffff"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.history_limit_days", "90"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.head_code", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.restricted_ips", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.support_email", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.locked_when_maintenance", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.custom_footer", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.custom_incident_types_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.slack_subscriptions_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.date_format", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.maintenance_notification_hours", "6"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.subdomain", "terraform-test-updated"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.twitter_public_screen_name", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.header_logo_text", "Test Status Page from Terraform"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.member_restricted", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.url", "terraform.test"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.status_ok_color", "48CBA5"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.uptime_graph_days", "90"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.subscribers_enabled", "true"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.display_about", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.%", "2"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.en.header_logo_text", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.en.public_company_name", "Public company name EN"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.fr.header_logo_text", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.fr.public_company_name", "Public company name FR"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.tweet_by_default", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.display_calendar", "true"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.email_templates_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.google_calendar_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.link_color", "0c91c3"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.email_layout_template", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.name", "Edited Test Status Page from Terraform"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.status_major_color", "e75a53"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.custom_header", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.date_format_enforce_everywhere", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.time_format", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.header_bg_color1", "009688"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.incident_link_color", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.bg_image", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.logo", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.favicon", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.custom_css", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.current_incidents_position", "below_services"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.custom_js", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.minor_notification_hours", "6"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.mattermost_notifications_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.info_notices_enabled", "true"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.captcha_enabled", "true"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.about", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.google_chat_notifications_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.discord_notifications_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.theme_configs.%", "9"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.theme_configs.link_color", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.theme_configs.header_bg_color1", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.theme_configs.header_bg_color2", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.theme_configs.header_fg_color", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.theme_configs.incident_header_color", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.theme_configs.status_ok_color", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.theme_configs.status_minor_color", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.theme_configs.status_major_color", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.theme_configs.status_maintenance_color", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.status_minor_color", "FFA500"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.tweeting_enabled", "true"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.sms_notifications_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.notify_by_default", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.hide_watermark", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.custom_domain_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.enable_auto_translations", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.feed_enabled", "true"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.header_bg_color2", "0c91c3"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.public_company_name", "Public company name EN"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.notification_email", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.email_notification_template", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.teams_notifications_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.status_maintenance_color", "5378c1"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.email_confirmation_template", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.time_zone", "Europe/Budapest"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.calendar_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.major_notification_hours", "3"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.incident_header_color", "009688"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.reply_to_email", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.noindex", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.inserted_at", "2024-04-15T11:20:35"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.updated_at", "2024-04-25T11:22:32"),
					// Verify placeholder id attribute
					resource.TestCheckResourceAttr("statuspal_status_page.test", "id", "placeholder"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
