package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
			"header_logo_text": "Test Status Page from Terraform EN",
			"member_restricted": false,
			"url": "terraform.test",
			"status_ok_color": "48CBA5",
			"uptime_graph_days": 90,
			"subscribers_enabled": true,
			"display_about": false,
			"translations": {
				"en": {
					"header_logo_text": "Test Status Page from Terraform EN",
					"public_company_name": "Public company name EN"
				},
				"fr": {
					"header_logo_text": "Test Status Page from Terraform FR",
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
			"status_minor_color": "FFA500",
			"tweeting_enabled": true,
			"sms_notifications_enabled": false,
			"zoom_notifications_enabled": false,
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
			"noindex": false,
			"allowed_email_domains": "acme.corp\nbbc.com",
			"domain_config": {
				"provider": "cloudflare",
				"domain": "status.terraform.test",
				"main_hostname": "ssl-for-saas.example.com",
				"status": "configuring",
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
								header_logo_text = "Test Status Page from Terraform EN"
								public_company_name = "Public company name EN"
							}
							fr = {
								header_logo_text = "Test Status Page from Terraform FR"
								public_company_name = "Public company name FR"
							}
						}
						allowed_email_domains = "acme.corp\nbbc.com"
						public_company_name = "Public company name EN"
						domain_config = {
							provider = "cloudflare"
							domain   = "status.terraform.test"
						}
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspal_status_page.test", "organization_id", "1"),
					// Verify status_page
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.%", "77"),
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
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.header_logo_text", "Test Status Page from Terraform EN"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.member_restricted", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.url", "terraform.test"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.status_ok_color", "48CBA5"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.uptime_graph_days", "90"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.subscribers_enabled", "true"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.display_about", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.%", "2"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.en.%", "2"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.en.header_logo_text", "Test Status Page from Terraform EN"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.en.public_company_name", "Public company name EN"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.fr.%", "2"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.fr.header_logo_text", "Test Status Page from Terraform FR"),
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
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.zoom_notifications_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.allowed_email_domains", "acme.corp\nbbc.com"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.inserted_at", "2024-04-15T11:20:35"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.updated_at", "2024-04-20T11:22:32"),
					// Verify domain_config
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.%", "8"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.provider", "cloudflare"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.domain", "status.terraform.test"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.main_hostname", "ssl-for-saas.example.com"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.status", "configuring"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.error", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.external_id", "ext-abc123"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.validation_records.%", "2"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.validation_records.cname.name", "status.terraform.test"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.validation_records.cname.type", "CNAME"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.validation_records.cname.value", "ssl-for-saas.example.com"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.validation_records.hostname_txt.name", "_cf-custom-hostname.status.terraform.test"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.validation_records.hostname_txt.type", "TXT"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.validation_records.hostname_txt.value", "some-verification-token"),
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
								header_logo_text = "Test Status Page from Terraform EN"
								public_company_name = "Public company name EN"
							}
							fr = {
								header_logo_text = "Test Status Page from Terraform FR"
								public_company_name = "Public company name FR"
							}
						}
						allowed_email_domains = "acme.corp\nbbc.com"
						public_company_name = "Public company name EN"
						domain_config = {
							provider = "cloudflare"
							domain   = "status.terraform.test"
						}
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspal_status_page.test", "organization_id", "1"),
					// Verify status_page
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.%", "77"),
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
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.header_logo_text", "Test Status Page from Terraform EN"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.member_restricted", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.url", "terraform.test"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.status_ok_color", "48CBA5"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.uptime_graph_days", "90"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.subscribers_enabled", "true"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.display_about", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.%", "2"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.en.%", "2"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.en.header_logo_text", "Test Status Page from Terraform EN"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.en.public_company_name", "Public company name EN"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.fr.%", "2"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.translations.fr.header_logo_text", "Test Status Page from Terraform FR"),
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
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.zoom_notifications_enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.allowed_email_domains", "acme.corp\nbbc.com"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.inserted_at", "2024-04-15T11:20:35"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.updated_at", "2024-04-25T11:22:32"),
					// Verify domain_config
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.%", "8"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.provider", "cloudflare"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.domain", "status.terraform.test"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.main_hostname", "ssl-for-saas.example.com"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.status", "configuring"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.error", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.external_id", "ext-abc123"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.validation_records.%", "2"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.validation_records.cname.name", "status.terraform.test"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.validation_records.cname.type", "CNAME"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.validation_records.cname.value", "ssl-for-saas.example.com"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.validation_records.hostname_txt.name", "_cf-custom-hostname.status.terraform.test"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.validation_records.hostname_txt.type", "TXT"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.validation_records.hostname_txt.value", "some-verification-token"),
					// Verify placeholder id attribute
					resource.TestCheckResourceAttr("statuspal_status_page.test", "id", "placeholder"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// TestAccStatusPageResource_LegacyDomain verifies that creating a status page with
// the deprecated domain + custom_domain_enabled fields works when the API converts
// them into domain_config with provider "legacy_custom_domain" and returns domain "".
func TestAccStatusPageResource_LegacyDomain(t *testing.T) {
	mux := http.NewServeMux()

	legacyResponse := `{
		"status_page": {
			"name": "Legacy Domain Test",
			"url": "legacy.test",
			"time_zone": "UTC",
			"subdomain": "legacy-test",
			"domain": "",
			"custom_domain_enabled": true,
			"domain_config": {
				"provider": "legacy_custom_domain",
				"domain": "status.legacy.test",
				"main_hostname": null,
				"status": "active",
				"error": null,
				"external_id": null,
				"pullzone_id": null,
				"validation_records": {}
			},
			"theme_selected": "default",
			"scheduled_maintenance_days": 7,
			"display_uptime_graph": true,
			"inserted_at": "2024-06-01T10:00:00",
			"updated_at": "2024-06-01T10:00:00",
			"header_fg_color": "ffffff",
			"history_limit_days": 90,
			"head_code": null,
			"support_email": null,
			"locked_when_maintenance": false,
			"custom_footer": null,
			"custom_incident_types_enabled": false,
			"slack_subscriptions_enabled": false,
			"date_format": null,
			"maintenance_notification_hours": 6,
			"twitter_public_screen_name": null,
			"header_logo_text": null,
			"member_restricted": false,
			"status_ok_color": "48CBA5",
			"uptime_graph_days": 90,
			"subscribers_enabled": true,
			"display_about": false,
			"translations": {},
			"tweet_by_default": false,
			"display_calendar": true,
			"email_templates_enabled": false,
			"google_calendar_enabled": false,
			"link_color": "0c91c3",
			"email_layout_template": null,
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
			"status_minor_color": "FFA500",
			"tweeting_enabled": true,
			"sms_notifications_enabled": false,
			"zoom_notifications_enabled": false,
			"notify_by_default": false,
			"hide_watermark": false,
			"enable_auto_translations": false,
			"restricted_ips": null,
			"feed_enabled": true,
			"header_bg_color2": "0c91c3",
			"public_company_name": null,
			"notification_email": null,
			"email_notification_template": null,
			"teams_notifications_enabled": false,
			"status_maintenance_color": "5378c1",
			"email_confirmation_template": null,
			"calendar_enabled": false,
			"major_notification_hours": 3,
			"incident_header_color": "009688",
			"reply_to_email": null,
			"noindex": false,
			"allowed_email_domains": null
		}
	}`

	mux.HandleFunc("/orgs/1/status_pages", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(legacyResponse)); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})
	mux.HandleFunc("/orgs/1/status_pages/legacy-test", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			w.Write([]byte(`""`))
			return
		}
		w.Write([]byte(legacyResponse))
	})

	mockServer := httptest.NewServer(mux)
	defer mockServer.Close()
	providerConfig := providerConfig(&mockServer.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: *providerConfig + `resource "statuspal_status_page" "test" {
					organization_id = "1"
					status_page = {
						name      = "Legacy Domain Test"
						url       = "legacy.test"
						time_zone = "UTC"
						domain                = "status.legacy.test"
						custom_domain_enabled = true
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// The provider should echo domain_config.domain back into the legacy domain field
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain", "status.legacy.test"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.custom_domain_enabled", "true"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.provider", "legacy_custom_domain"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.domain", "status.legacy.test"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.status", "active"),
				),
			},
		},
	})
}

// TestAccStatusPageResource_LegacyToCloudFlareMigration verifies that migrating from
// a legacy_custom_domain to cloudflare sends a clearing API call before the real update.
func TestAccStatusPageResource_LegacyToCloudFlareMigration(t *testing.T) {
	var clearCallCount atomic.Int32

	mux := http.NewServeMux()

	legacyResponse := `{
		"status_page": {
			"name": "Migration Test",
			"url": "migrate.test",
			"time_zone": "UTC",
			"subdomain": "migrate-test",
			"domain": "",
			"custom_domain_enabled": true,
			"domain_config": {
				"provider": "legacy_custom_domain",
				"domain": "status.migrate.test",
				"main_hostname": null,
				"status": "active",
				"error": null,
				"external_id": null,
				"pullzone_id": null,
				"validation_records": {}
			},
			"theme_selected": "default",
			"scheduled_maintenance_days": 7,
			"display_uptime_graph": true,
			"inserted_at": "2024-06-01T10:00:00",
			"updated_at": "2024-06-01T10:00:00",
			"header_fg_color": "ffffff",
			"history_limit_days": 90,
			"head_code": null,
			"support_email": null,
			"locked_when_maintenance": false,
			"custom_footer": null,
			"custom_incident_types_enabled": false,
			"slack_subscriptions_enabled": false,
			"date_format": null,
			"maintenance_notification_hours": 6,
			"twitter_public_screen_name": null,
			"header_logo_text": null,
			"member_restricted": false,
			"status_ok_color": "48CBA5",
			"uptime_graph_days": 90,
			"subscribers_enabled": true,
			"display_about": false,
			"translations": {},
			"tweet_by_default": false,
			"display_calendar": true,
			"email_templates_enabled": false,
			"google_calendar_enabled": false,
			"link_color": "0c91c3",
			"email_layout_template": null,
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
			"status_minor_color": "FFA500",
			"tweeting_enabled": true,
			"sms_notifications_enabled": false,
			"zoom_notifications_enabled": false,
			"notify_by_default": false,
			"hide_watermark": false,
			"enable_auto_translations": false,
			"restricted_ips": null,
			"feed_enabled": true,
			"header_bg_color2": "0c91c3",
			"public_company_name": null,
			"notification_email": null,
			"email_notification_template": null,
			"teams_notifications_enabled": false,
			"status_maintenance_color": "5378c1",
			"email_confirmation_template": null,
			"calendar_enabled": false,
			"major_notification_hours": 3,
			"incident_header_color": "009688",
			"reply_to_email": null,
			"noindex": false,
			"allowed_email_domains": null
		}
	}`

	cloudflareResponse := strings.Replace(legacyResponse,
		`"provider": "legacy_custom_domain"`,
		`"provider": "cloudflare"`, 1)
	cloudflareResponse = strings.Replace(cloudflareResponse,
		`"status": "active"`,
		`"status": "configuring"`, 1)
	cloudflareResponse = strings.Replace(cloudflareResponse,
		`"validation_records": {}`,
		`"validation_records": {
			"hostname_cname_name": "status.migrate.test",
			"hostname_cname_value": "domains-proxied.statuspal.io"
		}`, 1)
	cloudflareResponse = strings.Replace(cloudflareResponse,
		`"custom_domain_enabled": true`,
		`"custom_domain_enabled": false`, 1)

	currentResponse := legacyResponse

	mux.HandleFunc("/orgs/1/status_pages", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(currentResponse))
	})
	mux.HandleFunc("/orgs/1/status_pages/migrate-test", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			w.Write([]byte(`""`))
			return
		}
		if r.Method == http.MethodPut {
			body, _ := io.ReadAll(r.Body)
			var payload struct {
				StatusPage struct {
					DomainConfig json.RawMessage `json:"domain_config"`
				} `json:"status_page"`
			}
			json.Unmarshal(body, &payload)

			// Detect the clearing call: domain_config is absent (null/omitted)
			if payload.StatusPage.DomainConfig == nil || string(payload.StatusPage.DomainConfig) == "null" {
				clearCallCount.Add(1)
				// After clearing, return the legacy response (clear acknowledged)
				w.Write([]byte(legacyResponse))
				return
			}
			// The real update with cloudflare provider
			currentResponse = cloudflareResponse
		}
		w.Write([]byte(currentResponse))
	})

	mockServer := httptest.NewServer(mux)
	defer mockServer.Close()
	providerConfig := providerConfig(&mockServer.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with legacy domain
			{
				Config: *providerConfig + `resource "statuspal_status_page" "test" {
					organization_id = "1"
					status_page = {
						name      = "Migration Test"
						url       = "migrate.test"
						time_zone = "UTC"
						domain                = "status.migrate.test"
						custom_domain_enabled = true
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.provider", "legacy_custom_domain"),
				),
			},
			// Step 2: Migrate to cloudflare — the provider should send a clearing call first
			{
				Config: *providerConfig + `resource "statuspal_status_page" "test" {
					organization_id = "1"
					status_page = {
						name      = "Migration Test"
						url       = "migrate.test"
						time_zone = "UTC"
						domain_config = {
							provider = "cloudflare"
							domain   = "status.migrate.test"
						}
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain", ""),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.provider", "cloudflare"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.domain", "status.migrate.test"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.status", "configuring"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.validation_records.cname.name", "status.migrate.test"),
					resource.TestCheckResourceAttr("statuspal_status_page.test", "status_page.domain_config.validation_records.cname.value", "domains-proxied.statuspal.io"),
				),
			},
		},
		// Verify the clearing call was made during the migration step
		CheckDestroy: func(s *terraform.State) error {
			if clearCallCount.Load() == 0 {
				return fmt.Errorf("expected a clearing API call during legacy-to-cloudflare migration, but none was made")
			}
			return nil
		},
	})
}
