package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	statuspal "terraform-provider-statuspal/internal/client"
)

var validationRecordAttrTypes = map[string]attr.Type{
	"name":  types.StringType,
	"type":  types.StringType,
	"value": types.StringType,
}

var domainConfigAttrTypes = map[string]attr.Type{
	"provider":           types.StringType,
	"domain":             types.StringType,
	"main_hostname":      types.StringType,
	"validation_records": types.MapType{ElemType: types.ObjectType{AttrTypes: validationRecordAttrTypes}},
	"external_id":        types.StringType,
	"status":             types.StringType,
	"error":              types.StringType,
	"pullzone_id":        types.Int64Type,
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &statusPageResource{}
	_ resource.ResourceWithConfigure   = &statusPageResource{}
	_ resource.ResourceWithImportState = &statusPageResource{}
)

// NewStatusPageResource is a helper function to simplify the provider implementation.
func NewStatusPageResource() resource.Resource {
	return &statusPageResource{}
}

// statusPageResource is the resource implementation.
type statusPageResource struct {
	client *statuspal.Client
}

// statusPageResourceModel maps the resource schema data.
type statusPageResourceModel struct {
	ID             types.String    `tfsdk:"id"` // only for test case
	OrganizationID types.String    `tfsdk:"organization_id"`
	StatusPage     statusPageModel `tfsdk:"status_page"`
}

// domainConfigModel maps domain_config schema data.
type domainConfigModel struct {
	CDNProvider       types.String `tfsdk:"provider"`
	Domain            types.String `tfsdk:"domain"`
	MainHostname      types.String `tfsdk:"main_hostname"`
	ValidationRecords types.Map    `tfsdk:"validation_records"`
	ExternalID        types.String `tfsdk:"external_id"`
	Status            types.String `tfsdk:"status"`
	Error             types.String `tfsdk:"error"`
	PullzoneID        types.Int64  `tfsdk:"pullzone_id"`
}

// statusPageModel maps status_page schema data.
type statusPageModel struct {
	Name                           types.String `tfsdk:"name"`
	Url                            types.String `tfsdk:"url"`
	TimeZone                       types.String `tfsdk:"time_zone"`
	Subdomain                      types.String `tfsdk:"subdomain"`
	SupportEmail                   types.String `tfsdk:"support_email"`
	TwitterPublicScreenName        types.String `tfsdk:"twitter_public_screen_name"`
	About                          types.String `tfsdk:"about"`
	DisplayAbout                   types.Bool   `tfsdk:"display_about"`
	CustomDomainEnabled            types.Bool   `tfsdk:"custom_domain_enabled"`
	Domain                         types.String `tfsdk:"domain"`
	DomainConfig                   types.Object `tfsdk:"domain_config"`
	RestrictedIps                  types.String `tfsdk:"restricted_ips"`
	MemberRestricted               types.Bool   `tfsdk:"member_restricted"`
	ScheduledMaintenanceDays       types.Int64  `tfsdk:"scheduled_maintenance_days"`
	CustomJs                       types.String `tfsdk:"custom_js"`
	HeadCode                       types.String `tfsdk:"head_code"`
	DateFormat                     types.String `tfsdk:"date_format"`
	TimeFormat                     types.String `tfsdk:"time_format"`
	DateFormatEnforceEverywhere    types.Bool   `tfsdk:"date_format_enforce_everywhere"`
	DisplayCalendar                types.Bool   `tfsdk:"display_calendar"`
	HideWatermark                  types.Bool   `tfsdk:"hide_watermark"`
	MinorNotificationHours         types.Int64  `tfsdk:"minor_notification_hours"`
	MajorNotificationHours         types.Int64  `tfsdk:"major_notification_hours"`
	MaintenanceNotificationHours   types.Int64  `tfsdk:"maintenance_notification_hours"`
	HistoryLimitDays               types.Int64  `tfsdk:"history_limit_days"`
	CustomIncidentTypesEnabled     types.Bool   `tfsdk:"custom_incident_types_enabled"`
	InfoNoticesEnabled             types.Bool   `tfsdk:"info_notices_enabled"`
	LockedWhenMaintenance          types.Bool   `tfsdk:"locked_when_maintenance"`
	Noindex                        types.Bool   `tfsdk:"noindex"`
	EnableAutoTranslations         types.Bool   `tfsdk:"enable_auto_translations"`
	CaptchaEnabled                 types.Bool   `tfsdk:"captcha_enabled"`
	Translations                   types.Map    `tfsdk:"translations"`
	HeaderLogoText                 types.String `tfsdk:"header_logo_text"`
	PublicCompanyName              types.String `tfsdk:"public_company_name"`
	BgImage                        types.String `tfsdk:"bg_image"`
	Logo                           types.String `tfsdk:"logo"`
	Favicon                        types.String `tfsdk:"favicon"`
	DisplayUptimeGraph             types.Bool   `tfsdk:"display_uptime_graph"`
	UptimeGraphDays                types.Int64  `tfsdk:"uptime_graph_days"`
	CurrentIncidentsPosition       types.String `tfsdk:"current_incidents_position"`
	ThemeSelected                  types.String `tfsdk:"theme_selected"`
	LinkColor                      types.String `tfsdk:"link_color"`
	HeaderBgColor1                 types.String `tfsdk:"header_bg_color1"`
	HeaderBgColor2                 types.String `tfsdk:"header_bg_color2"`
	HeaderFgColor                  types.String `tfsdk:"header_fg_color"`
	IncidentHeaderColor            types.String `tfsdk:"incident_header_color"`
	IncidentLinkColor              types.String `tfsdk:"incident_link_color"`
	StatusOkColor                  types.String `tfsdk:"status_ok_color"`
	StatusMinorColor               types.String `tfsdk:"status_minor_color"`
	StatusMajorColor               types.String `tfsdk:"status_major_color"`
	StatusMaintenanceColor         types.String `tfsdk:"status_maintenance_color"`
	CustomCss                      types.String `tfsdk:"custom_css"`
	CustomHeader                   types.String `tfsdk:"custom_header"`
	CustomFooter                   types.String `tfsdk:"custom_footer"`
	NotifyByDefault                types.Bool   `tfsdk:"notify_by_default"`
	TweetByDefault                 types.Bool   `tfsdk:"tweet_by_default"`
	SlackSubscriptionsEnabled      types.Bool   `tfsdk:"slack_subscriptions_enabled"`
	DiscordNotificationsEnabled    types.Bool   `tfsdk:"discord_notifications_enabled"`
	TeamsNotificationsEnabled      types.Bool   `tfsdk:"teams_notifications_enabled"`
	GoogleChatNotificationsEnabled types.Bool   `tfsdk:"google_chat_notifications_enabled"`
	MattermostNotificationsEnabled types.Bool   `tfsdk:"mattermost_notifications_enabled"`
	SmsNotificationsEnabled        types.Bool   `tfsdk:"sms_notifications_enabled"`
	ZoomNotificationsEnabled       types.Bool   `tfsdk:"zoom_notifications_enabled"`
	FeedEnabled                    types.Bool   `tfsdk:"feed_enabled"`
	CalendarEnabled                types.Bool   `tfsdk:"calendar_enabled"`
	GoogleCalendarEnabled          types.Bool   `tfsdk:"google_calendar_enabled"`
	SubscribersEnabled             types.Bool   `tfsdk:"subscribers_enabled"`
	NotificationEmail              types.String `tfsdk:"notification_email"`
	ReplyToEmail                   types.String `tfsdk:"reply_to_email"`
	TweetingEnabled                types.Bool   `tfsdk:"tweeting_enabled"`
	EmailLayoutTemplate            types.String `tfsdk:"email_layout_template"`
	EmailConfirmationTemplate      types.String `tfsdk:"email_confirmation_template"`
	EmailNotificationTemplate      types.String `tfsdk:"email_notification_template"`
	EmailTemplatesEnabled          types.Bool   `tfsdk:"email_templates_enabled"`
	AllowedEmailDomains            types.String `tfsdk:"allowed_email_domains"`
	InsertedAt                     types.String `tfsdk:"inserted_at"`
	UpdatedAt                      types.String `tfsdk:"updated_at"`
}

type statusPageTranslationsModel map[string]statusPageTranslationModel

type statusPageTranslationModel struct {
	PublicCompanyName types.String `tfsdk:"public_company_name"`
	HeaderLogoText    types.String `tfsdk:"header_logo_text"`
}

// Metadata returns the resource type name.
func (r *statusPageResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_status_page"
}

// Schema defines the schema for the resource.
func (r *statusPageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a status page of the organization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier attribute. Ignore it, only used in testing.",
				Computed:    true,
			},
			"organization_id": schema.StringAttribute{
				Description: "The organization ID of the status page.",
				Required:    true,
			},
			"status_page": schema.SingleNestedAttribute{
				Description: "The status page.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "Company, project or service name.",
						Required:    true,
					},
					"url": schema.StringAttribute{
						Description: "The website to your company, project or service.",
						Required:    true,
					},
					"time_zone": schema.StringAttribute{
						Description: `The primary timezone the status page uses to display incidents (e.g. "Europe/Berlin").`,
						Required:    true,
					},
					"subdomain": schema.StringAttribute{
						Description: "The status page subdomain on statuspal.",
						Optional:    true,
						Computed:    true,
					},
					"support_email": schema.StringAttribute{
						Description: "Your company's support email.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"twitter_public_screen_name": schema.StringAttribute{
						Description: "Twitter handle name (e.g. yourcompany).",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"about": schema.StringAttribute{
						Description: "Customize the about information displayed in your status page.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"display_about": schema.BoolAttribute{
						Description: "Display about information.",
						Optional:    true,
						Computed:    true,
					},
					"custom_domain_enabled": schema.BoolAttribute{
						Description:        "Enable your custom domain with SSL.",
						DeprecationMessage: "Use domain_config instead.",
						Optional:           true,
						Computed:           true,
					},
					"domain": schema.StringAttribute{
						Description:        "Configure your own domain to point to your status page (e.g. status.your-company.com), we generate and auto-renew its SSL certificate for you.",
						DeprecationMessage: "Use domain_config.domain instead.",
						Optional:           true,
						Computed:           true,
						Default:            stringdefault.StaticString(""),
					},
					"restricted_ips": schema.StringAttribute{
						Description: `Your status page will be accessible only from this IPs (e.g. "1.1.1.1, 2.2.2.2").`,
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"member_restricted": schema.BoolAttribute{
						Description: "Only signed in members will be allowed to access your status page.",
						Optional:    true,
						Computed:    true,
					},
					"scheduled_maintenance_days": schema.Int64Attribute{
						Description: "Display scheduled maintenance.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(7),
						Validators: []validator.Int64{
							int64validator.OneOf(7, 14, 21, 28),
						},
					},
					"custom_js": schema.StringAttribute{
						MarkdownDescription: "We'll insert this content inside the `<script>` tag at the bottom of your status page `<body>` tag.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"head_code": schema.StringAttribute{
						MarkdownDescription: "We'll insert this content inside the `<head>` tag.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"date_format": schema.StringAttribute{
						Description: "Display timestamps of incidents and updates in this format.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"time_format": schema.StringAttribute{
						Description: "Display timestamps of incidents and updates in this format.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"date_format_enforce_everywhere": schema.BoolAttribute{
						Description: "The above date format will be used everywhere in the status page. Timezone conversion to client's will be disabled.",
						Optional:    true,
						Computed:    true,
					},
					"display_calendar": schema.BoolAttribute{
						Description: "Display uptime calendar at status page.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"hide_watermark": schema.BoolAttribute{
						Description: `Hide "Powered by Statuspal.io".`,
						Optional:    true,
						Computed:    true,
					},
					"minor_notification_hours": schema.Int64Attribute{
						Description: "Long-running incident notification (Minor incident).",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(6),
						Validators: []validator.Int64{
							int64validator.AtLeast(0),
						},
					},
					"major_notification_hours": schema.Int64Attribute{
						Description: "Long-running incident notification (Major incident).",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(3),
						Validators: []validator.Int64{
							int64validator.AtLeast(0),
						},
					},
					"maintenance_notification_hours": schema.Int64Attribute{
						Description: "Long-running incident notification (Maintenance).",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(6),
						Validators: []validator.Int64{
							int64validator.AtLeast(0),
						},
					},
					"history_limit_days": schema.Int64Attribute{
						Description: "Incident history limit (omit for No Limit).",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(90),
						Validators: []validator.Int64{
							int64validator.OneOf(30, 90, 365),
						},
					},
					"custom_incident_types_enabled": schema.BoolAttribute{
						Description: "Enable custom incident types.",
						Optional:    true,
						Computed:    true,
					},
					"info_notices_enabled": schema.BoolAttribute{
						Description: "Enable information notices.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"locked_when_maintenance": schema.BoolAttribute{
						Description: "Lock from adding incidents when under maintenance.",
						Optional:    true,
						Computed:    true,
					},
					"noindex": schema.BoolAttribute{
						Description: "Remove status page from being indexed by search engines (e.g. Google).",
						Optional:    true,
						Computed:    true,
					},
					"enable_auto_translations": schema.BoolAttribute{
						Description: "Enable auto translations when creating incidents, maintenances and info notices.",
						Optional:    true,
						Computed:    true,
					},
					"captcha_enabled": schema.BoolAttribute{
						Description: "Enable captchas (this option is only available when the status page is member restricted).",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"translations": schema.MapNestedAttribute{
						MarkdownDescription: "A translations object. For example:\n  ```terraform" + `
	{
		en = {
			public_company_name = "Your company"
			header_logo_text = "Your company status page"
		}
		fr = {
			public_company_name = "Votre entreprise"
			header_logo_text = "Page d'état de votre entreprise"
		}
	}
` + "  ```\n→ ",
						Optional: true,
						Computed: true,
						Default:  mapdefault.StaticValue(types.MapNull(types.ObjectType{AttrTypes: map[string]attr.Type{"public_company_name": types.StringType, "header_logo_text": types.StringType}})),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"public_company_name": schema.StringAttribute{
									Description: "Displayed at the footer of the status page.",
									Required:    true,
								},
								"header_logo_text": schema.StringAttribute{
									Description: "Displayed at the header of the status page.",
									Required:    true,
								},
							},
						},
					},
					"header_logo_text": schema.StringAttribute{
						Description: "Displayed at the header of the status page.",
						Optional:    true,
						Computed:    true,
					},
					"public_company_name": schema.StringAttribute{
						Description: "Displayed at the footer of the status page.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"bg_image": schema.StringAttribute{
						Description: "Background image url of the status page.",
						Computed:    true,
					},
					"logo": schema.StringAttribute{
						Description: "Logo url of the status page.",
						Computed:    true,
					},
					"favicon": schema.StringAttribute{
						Description: "Favicon url of the status page.",
						Computed:    true,
					},
					"display_uptime_graph": schema.BoolAttribute{
						Description: "Display the uptime graph in the status page.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"uptime_graph_days": schema.Int64Attribute{
						Description: "Uptime graph period.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(90),
						Validators: []validator.Int64{
							int64validator.OneOf(30, 60, 90),
						},
					},
					"current_incidents_position": schema.StringAttribute{
						Description: `The incident position displayed in the status page, it can be "below_services" and "above_services".`,
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("below_services"),
						Validators: []validator.String{
							stringvalidator.OneOf("below_services", "above_services"),
						},
					},
					"theme_selected": schema.StringAttribute{
						Description: `The selected theme for state page, it can be "default" and "big-logo".`,
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("default"),
						Validators: []validator.String{
							stringvalidator.OneOf("default", "big-logo"),
						},
					},
					"link_color": schema.StringAttribute{
						Description: "The links color in the status page.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("0c91c3"),
					},
					"header_bg_color1": schema.StringAttribute{
						Description: "The background color at left side of the status page header.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("009688"),
					},
					"header_bg_color2": schema.StringAttribute{
						Description: "The background color at right side of the status page header.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("0c91c3"),
					},
					"header_fg_color": schema.StringAttribute{
						Description: "The text color in the status page.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("ffffff"),
					},
					"incident_header_color": schema.StringAttribute{
						Description: "Incidents header color in the status page.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("009688"),
					},
					"incident_link_color": schema.StringAttribute{
						Description: "Incidents link color in the status page.",
						Optional:    true,
						Computed:    true,
					},
					"status_ok_color": schema.StringAttribute{
						Description: "The status page colors when there is no incident.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("48CBA5"),
					},
					"status_minor_color": schema.StringAttribute{
						Description: "The status page colors when there is a minor incident.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("FFA500"),
					},
					"status_major_color": schema.StringAttribute{
						Description: "The status page colors when there is a major incident.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("e75a53"),
					},
					"status_maintenance_color": schema.StringAttribute{
						Description: "The status page colors when there is a maintenance incident.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("5378c1"),
					},
					"custom_css": schema.StringAttribute{
						MarkdownDescription: "We'll insert this content inside the `<style>` tag.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"custom_header": schema.StringAttribute{
						MarkdownDescription: "A custom header for the status page (e.g. \"`<header>...</header>`\").",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"custom_footer": schema.StringAttribute{
						MarkdownDescription: "A custom footer for the status page (e.g. \"`<footer>...</footer>`\").",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"notify_by_default": schema.BoolAttribute{
						Description: "Check the Notify subscribers checkbox by default.",
						Optional:    true,
						Computed:    true,
					},
					"tweet_by_default": schema.BoolAttribute{
						Description: "Check the Tweet checkbox by default.",
						Optional:    true,
						Computed:    true,
					},
					"slack_subscriptions_enabled": schema.BoolAttribute{
						Description: "Allow your customers to subscribe via Slack to updates on your status page's status.",
						Optional:    true,
						Computed:    true,
					},
					"discord_notifications_enabled": schema.BoolAttribute{
						Description: "Allow your customers to receive notifications on a Discord channel.",
						Optional:    true,
						Computed:    true,
					},
					"teams_notifications_enabled": schema.BoolAttribute{
						Description: "Allow your customers to receive notifications on Microsoft Teams.",
						Optional:    true,
						Computed:    true,
					},
					"google_chat_notifications_enabled": schema.BoolAttribute{
						Description: "Allow your customers to receive notifications on Google Chat.",
						Optional:    true,
						Computed:    true,
					},
					"mattermost_notifications_enabled": schema.BoolAttribute{
						Description: "Allow your customers to receive notifications on Mattermost.",
						Optional:    true,
						Computed:    true,
					},
					"sms_notifications_enabled": schema.BoolAttribute{
						Description: "Allow your customers to receive SMS notifications on your status page's status (to enable this you need to have a Twilio or Esendex integration).",
						Optional:    true,
						Computed:    true,
					},
					"zoom_notifications_enabled": schema.BoolAttribute{
						Description: "Allow your customers to receive notifications on Zoom.",
						Optional:    true,
						Computed:    true,
					},
					"feed_enabled": schema.BoolAttribute{
						Description: "Allow your customers to receive updates as RSS and Atom feeds.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"calendar_enabled": schema.BoolAttribute{
						Description: "Allow your customers to receive updates via iCalendar feed.",
						Optional:    true,
						Computed:    true,
					},
					"google_calendar_enabled": schema.BoolAttribute{
						Description: "Allow your customers to import Google Calendar with Status Pages maintenance (business only).",
						Optional:    true,
						Computed:    true,
					},
					"subscribers_enabled": schema.BoolAttribute{
						Description: "Allow email customers to receive email notifications.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"notification_email": schema.StringAttribute{
						Description: "Allow your customers to subscribe via email to updates on your status page's status.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"reply_to_email": schema.StringAttribute{
						Description: "The email address we'll use in the 'reply_to' field in emails to your subscribers. So they can reply to your notification emails.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"tweeting_enabled": schema.BoolAttribute{
						Description: "Allows to send tweets when creating or updating an incident.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"email_layout_template": schema.StringAttribute{
						MarkdownDescription: "Custom email layout template, see the documentation: [Custom email templates](https://docs.statuspal.io/platform/subscriptions-and-notifications/custom-email-templates).",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"email_confirmation_template": schema.StringAttribute{
						MarkdownDescription: "Custom confirmation email template, see the documentation: [Custom email templates](https://docs.statuspal.io/platform/subscriptions-and-notifications/custom-email-templates).",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"email_notification_template": schema.StringAttribute{
						MarkdownDescription: "Custom email notification template, see the documentation: [Custom email templates](https://docs.statuspal.io/platform/subscriptions-and-notifications/custom-email-templates).",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"email_templates_enabled": schema.BoolAttribute{
						Description: "The templates won't be used until this is enabled, but you can send test emails.",
						Optional:    true,
						Computed:    true,
					},
					"allowed_email_domains": schema.StringAttribute{
						MarkdownDescription: "Users with these domains in their email address will be able to sign up via status page invite link. Each domain should be separated by `\\n` (e.g., `acme.corp\\nnapster.com`).",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"inserted_at": schema.StringAttribute{
						Description: "Datetime at which the status page was inserted.",
						Computed:    true,
					},
					"updated_at": schema.StringAttribute{
						Description: "Datetime at which the status page was last updated.",
						Computed:    true,
					},
					"domain_config": schema.SingleNestedAttribute{
						Description: "Custom domain configuration for the status page.",
						Optional:    true,
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"provider": schema.StringAttribute{
								Description: `Custom domain provider, either "cloudflare" or "bunny".`,
								Optional:    true,
								Computed:    true,
								Validators: []validator.String{
									stringvalidator.OneOf("cloudflare", "bunny"),
								},
							},
							"domain": schema.StringAttribute{
								Description: `The custom hostname (e.g. "status.acme.com"). Must be lowercase.`,
								Optional:    true,
								Computed:    true,
								Validators: []validator.String{
									stringvalidator.RegexMatches(
										regexp.MustCompile(`^[^A-Z]*$`),
										"must be lowercase",
									),
								},
							},
							"main_hostname": schema.StringAttribute{
								Description: "The CNAME target to point your domain at.",
								Computed:    true,
							},
							"validation_records": schema.MapNestedAttribute{
								Description: `DNS records required for domain setup. Keys: "cname" (CNAME routing record), "hostname_txt" (TXT record for hostname verification), "txt" (TXT record for ACME SSL challenge — only present after CNAME is in DNS). Not all keys are present at every lifecycle stage.`,
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"name": schema.StringAttribute{
											Description: "DNS record name (hostname).",
											Computed:    true,
										},
										"type": schema.StringAttribute{
											Description: `DNS record type (e.g. "CNAME", "TXT").`,
											Computed:    true,
										},
										"value": schema.StringAttribute{
											Description: "DNS record value.",
											Computed:    true,
										},
									},
								},
							},
							"external_id": schema.StringAttribute{
								Description: "Upstream provider identifier, useful for debugging.",
								Computed:    true,
							},
							"status": schema.StringAttribute{
								Description: `Current verification state: "disabled", "configuring", "active", or "failed_to_configure".`,
								Computed:    true,
							},
							"error": schema.StringAttribute{
								Description: `Error details when status is "failed_to_configure".`,
								Computed:    true,
							},
							"pullzone_id": schema.Int64Attribute{
								Description: "Bunny-specific pullzone ID.",
								Computed:    true,
							},
						},
					},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *statusPageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan statusPageResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	statusPage := mapStatusPageModelToRequestBody(&ctx, &plan.StatusPage, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new status page
	organizationID := plan.OrganizationID.ValueString()
	newStatusPage, err := r.client.CreateStatusPage(statusPage, &organizationID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating StatusPal StatusPage",
			"Could not create status page, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	newStatusPageModel := mapResponseToStatusPageModel(newStatusPage, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.StatusPage = *newStatusPageModel
	plan.ID = types.StringValue("placeholder") // only for test case

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *statusPageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state statusPageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed status page value from StatusPal
	organizationID := state.OrganizationID.ValueString()
	subdomain := state.StatusPage.Subdomain.ValueString()
	statusPage, err := r.client.GetStatusPage(&organizationID, &subdomain)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading StatusPal StatusPage",
			"Could not read status page subdomain "+subdomain+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	statusPageModel := mapResponseToStatusPageModel(statusPage, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	state.StatusPage = *statusPageModel
	state.ID = types.StringValue("placeholder") // only for test case

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *statusPageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state statusPageResourceModel
	stateDiags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from plan
	var plan statusPageResourceModel
	planDiags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(planDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	statusPage := mapStatusPageModelToRequestBody(&ctx, &plan.StatusPage, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if statusPage.Subdomain == "" {
		statusPage.Subdomain = state.StatusPage.Subdomain.ValueString()
	}

	// Update existing status page
	organizationID := plan.OrganizationID.ValueString()
	subdomain := state.StatusPage.Subdomain.ValueString()
	updatedStatusPage, err := r.client.UpdateStatusPage(statusPage, &organizationID, &subdomain)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating StatusPal StatusPage",
			"Could not Update status page, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	updatedStatusPageModel := mapResponseToStatusPageModel(updatedStatusPage, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.StatusPage = *updatedStatusPageModel
	plan.ID = types.StringValue("placeholder") // only for test case

	// Set state to fully populated data
	stateDiags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *statusPageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state statusPageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	organizationID := state.OrganizationID.ValueString()
	subdomain := state.StatusPage.Subdomain.ValueString()
	err := r.client.DeleteStatusPage(&organizationID, &subdomain)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting StatusPal StatusPage",
			"Could not delete status page, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *statusPageResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Split the ID based on the delimiter used during import
	parts := strings.Split(req.ID, " ")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected StatusPal StatusPage Import Identifier",
			`Expected StatusPal status page import identifier with format: "<organization_id> <status_page_subdomain>"`,
		)
		return
	}

	req.ID = parts[0]
	resource.ImportStatePassthroughID(ctx, path.Root("organization_id"), req, resp)
	req.ID = parts[1]
	resource.ImportStatePassthroughID(ctx, path.Root("status_page").AtName("subdomain"), req, resp)
}

// Configure adds the provider configured client to the resource.
func (r *statusPageResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*statuspal.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf(
				"Expected *statuspal.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	r.client = client
}

func mapStatusPageModelToRequestBody(
	ctx *context.Context,
	statusPage *statusPageModel,
	diagnostics *diag.Diagnostics,
) *statuspal.StatusPage {
	// Create the translationData object dynamically
	translationData := make(statuspal.StatusPageTranslations)
	if !statusPage.Translations.IsNull() && !statusPage.Translations.IsUnknown() {
		translations := make(statusPageTranslationsModel, len(statusPage.Translations.Elements()))
		diags := statusPage.Translations.ElementsAs(*ctx, &translations, false)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return nil
		}

		for lang, data := range translations {
			translationData[lang] = statuspal.StatusPageTranslation{
				PublicCompanyName: data.PublicCompanyName.ValueString(),
				HeaderLogoText:    data.HeaderLogoText.ValueString(),
			}
		}
	}

	var domainConfig *statuspal.DomainConfig
	if !statusPage.DomainConfig.IsNull() && !statusPage.DomainConfig.IsUnknown() {
		var dc domainConfigModel
		diags := statusPage.DomainConfig.As(*ctx, &dc, basetypes.ObjectAsOptions{})
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return nil
		}
		provider := strings.ToLower(dc.CDNProvider.ValueString())
		domain := strings.ToLower(dc.Domain.ValueString())
		domainConfig = &statuspal.DomainConfig{
			CDNProvider: &provider,
			Domain:      &domain,
		}
	}

	return &statuspal.StatusPage{
		Name:                           statusPage.Name.ValueString(),
		Url:                            statusPage.Url.ValueString(),
		TimeZone:                       statusPage.TimeZone.ValueString(),
		Subdomain:                      statusPage.Subdomain.ValueString(),
		SupportEmail:                   statusPage.SupportEmail.ValueString(),
		TwitterPublicScreenName:        statusPage.TwitterPublicScreenName.ValueString(),
		About:                          statusPage.About.ValueString(),
		DisplayAbout:                   statusPage.DisplayAbout.ValueBool(),
		CustomDomainEnabled:            statusPage.CustomDomainEnabled.ValueBool(),
		Domain:                         statusPage.Domain.ValueString(),
		DomainConfig:                   domainConfig,
		RestrictedIps:                  statusPage.RestrictedIps.ValueString(),
		MemberRestricted:               statusPage.MemberRestricted.ValueBool(),
		ScheduledMaintenanceDays:       statusPage.ScheduledMaintenanceDays.ValueInt64(),
		CustomJs:                       statusPage.CustomJs.ValueString(),
		HeadCode:                       statusPage.HeadCode.ValueString(),
		DateFormat:                     statusPage.DateFormat.ValueString(),
		TimeFormat:                     statusPage.TimeFormat.ValueString(),
		DateFormatEnforceEverywhere:    statusPage.DateFormatEnforceEverywhere.ValueBool(),
		DisplayCalendar:                statusPage.DisplayCalendar.ValueBool(),
		HideWatermark:                  statusPage.HideWatermark.ValueBool(),
		MinorNotificationHours:         statusPage.MinorNotificationHours.ValueInt64(),
		MajorNotificationHours:         statusPage.MajorNotificationHours.ValueInt64(),
		MaintenanceNotificationHours:   statusPage.MaintenanceNotificationHours.ValueInt64(),
		HistoryLimitDays:               statusPage.HistoryLimitDays.ValueInt64(),
		CustomIncidentTypesEnabled:     statusPage.CustomIncidentTypesEnabled.ValueBool(),
		InfoNoticesEnabled:             statusPage.InfoNoticesEnabled.ValueBool(),
		LockedWhenMaintenance:          statusPage.LockedWhenMaintenance.ValueBool(),
		Noindex:                        statusPage.Noindex.ValueBool(),
		EnableAutoTranslations:         statusPage.EnableAutoTranslations.ValueBool(),
		CaptchaEnabled:                 statusPage.CaptchaEnabled.ValueBool(),
		Translations:                   translationData,
		HeaderLogoText:                 statusPage.HeaderLogoText.ValueString(),
		PublicCompanyName:              statusPage.PublicCompanyName.ValueString(),
		DisplayUptimeGraph:             statusPage.DisplayUptimeGraph.ValueBool(),
		UptimeGraphDays:                statusPage.UptimeGraphDays.ValueInt64(),
		CurrentIncidentsPosition:       statusPage.CurrentIncidentsPosition.ValueString(),
		ThemeSelected:                  statusPage.ThemeSelected.ValueString(),
		LinkColor:                      statusPage.LinkColor.ValueString(),
		HeaderBgColor1:                 statusPage.HeaderBgColor1.ValueString(),
		HeaderBgColor2:                 statusPage.HeaderBgColor2.ValueString(),
		HeaderFgColor:                  statusPage.HeaderFgColor.ValueString(),
		IncidentHeaderColor:            statusPage.IncidentHeaderColor.ValueString(),
		IncidentLinkColor:              statusPage.IncidentLinkColor.ValueString(),
		StatusOkColor:                  statusPage.StatusOkColor.ValueString(),
		StatusMinorColor:               statusPage.StatusMinorColor.ValueString(),
		StatusMajorColor:               statusPage.StatusMajorColor.ValueString(),
		StatusMaintenanceColor:         statusPage.StatusMaintenanceColor.ValueString(),
		CustomCss:                      statusPage.CustomCss.ValueString(),
		CustomHeader:                   statusPage.CustomHeader.ValueString(),
		CustomFooter:                   statusPage.CustomFooter.ValueString(),
		NotifyByDefault:                statusPage.NotifyByDefault.ValueBool(),
		TweetByDefault:                 statusPage.TweetByDefault.ValueBool(),
		SlackSubscriptionsEnabled:      statusPage.SlackSubscriptionsEnabled.ValueBool(),
		DiscordNotificationsEnabled:    statusPage.DiscordNotificationsEnabled.ValueBool(),
		TeamsNotificationsEnabled:      statusPage.TeamsNotificationsEnabled.ValueBool(),
		GoogleChatNotificationsEnabled: statusPage.GoogleChatNotificationsEnabled.ValueBool(),
		MattermostNotificationsEnabled: statusPage.MattermostNotificationsEnabled.ValueBool(),
		SmsNotificationsEnabled:        statusPage.SmsNotificationsEnabled.ValueBool(),
		FeedEnabled:                    statusPage.FeedEnabled.ValueBool(),
		CalendarEnabled:                statusPage.CalendarEnabled.ValueBool(),
		GoogleCalendarEnabled:          statusPage.GoogleCalendarEnabled.ValueBool(),
		SubscribersEnabled:             statusPage.SubscribersEnabled.ValueBool(),
		NotificationEmail:              statusPage.NotificationEmail.ValueString(),
		ReplyToEmail:                   statusPage.ReplyToEmail.ValueString(),
		TweetingEnabled:                statusPage.TweetingEnabled.ValueBool(),
		EmailLayoutTemplate:            statusPage.EmailLayoutTemplate.ValueString(),
		EmailConfirmationTemplate:      statusPage.EmailConfirmationTemplate.ValueString(),
		EmailNotificationTemplate:      statusPage.EmailNotificationTemplate.ValueString(),
		EmailTemplatesEnabled:          statusPage.EmailTemplatesEnabled.ValueBool(),
		ZoomNotificationsEnabled:       statusPage.ZoomNotificationsEnabled.ValueBool(),
		AllowedEmailDomains:            statusPage.AllowedEmailDomains.ValueString(),
	}
}

// buildValidationRecords converts the flat API map into a map of objects keyed by record type.
// Well-known prefixes are mapped to friendly keys and types:
//   - hostname_cname → "cname", type "CNAME"
//   - hostname_txt → "hostname_txt", type "TXT"
//   - certificate_txt → "txt", type "TXT"
//
// All other _name/_value pairs use the raw prefix as the key with an empty type.
func buildValidationRecords(raw map[string]string) (types.Map, diag.Diagnostics) {
	vrElemType := types.ObjectType{AttrTypes: validationRecordAttrTypes}
	wellKnown := map[string]struct{ key, recordType string }{
		"hostname_cname":  {"cname", "CNAME"},
		"hostname_txt":    {"hostname_txt", "TXT"},
		"certificate_txt": {"txt", "TXT"},
	}

	type record struct{ name, recordType, value string }
	collected := map[string]*record{}

	for rawKey, rawVal := range raw {
		var prefix, field string
		if strings.HasSuffix(rawKey, "_name") {
			prefix = strings.TrimSuffix(rawKey, "_name")
			field = "name"
		} else if strings.HasSuffix(rawKey, "_value") {
			prefix = strings.TrimSuffix(rawKey, "_value")
			field = "value"
		} else {
			continue
		}

		key, recordType := prefix, ""
		if meta, ok := wellKnown[prefix]; ok {
			key, recordType = meta.key, meta.recordType
		}

		if _, ok := collected[key]; !ok {
			collected[key] = &record{recordType: recordType}
		}
		if field == "name" {
			collected[key].name = rawVal
		} else {
			collected[key].value = rawVal
		}
	}

	elements := make(map[string]attr.Value, len(collected))
	var allDiags diag.Diagnostics
	for key, rec := range collected {
		if rec.name == "" {
			continue
		}
		obj, diags := types.ObjectValue(validationRecordAttrTypes, map[string]attr.Value{
			"name":  types.StringValue(rec.name),
			"type":  types.StringValue(rec.recordType),
			"value": types.StringValue(rec.value),
		})
		allDiags.Append(diags...)
		if allDiags.HasError() {
			return types.MapNull(vrElemType), allDiags
		}
		elements[key] = obj
	}

	result, diags := types.MapValue(vrElemType, elements)
	allDiags.Append(diags...)
	return result, allDiags
}

func mapResponseToStatusPageModel(statusPage *statuspal.StatusPage, diagnostics *diag.Diagnostics) *statusPageModel {
	// Map domain_config
	domainConfig := types.ObjectNull(domainConfigAttrTypes)
	if statusPage.DomainConfig != nil {
		dc := statusPage.DomainConfig

		vrElemType := types.ObjectType{AttrTypes: validationRecordAttrTypes}
		validationRecords := types.MapNull(vrElemType)
		if dc.ValidationRecords != nil {
			vr, diags := buildValidationRecords(dc.ValidationRecords)
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return nil
			}
			validationRecords = vr
		}

		var pullzoneID types.Int64
		if dc.PullzoneID != nil {
			pullzoneID = types.Int64Value(*dc.PullzoneID)
		} else {
			pullzoneID = types.Int64Null()
		}

		dcObj, diags := types.ObjectValue(domainConfigAttrTypes, map[string]attr.Value{
			"provider":           types.StringValue(stringPtrOrEmpty(dc.CDNProvider)),
			"domain":             types.StringValue(strings.ToLower(stringPtrOrEmpty(dc.Domain))),
			"main_hostname":      types.StringValue(stringPtrOrEmpty(dc.MainHostname)),
			"validation_records": validationRecords,
			"external_id":        types.StringValue(stringPtrOrEmpty(dc.ExternalID)),
			"status":             types.StringValue(stringPtrOrEmpty(dc.Status)),
			"error":              types.StringValue(stringPtrOrEmpty(dc.Error)),
			"pullzone_id":        pullzoneID,
		})
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return nil
		}
		domainConfig = dcObj
	}

	// Define the translation object schema
	translationSchema := map[string]attr.Type{
		"public_company_name": types.StringType,
		"header_logo_text":    types.StringType,
	}
	translations := types.MapNull(types.ObjectType{AttrTypes: translationSchema})
	if len(statusPage.Translations) > 0 {
		// Create the translationData object dynamically
		translationData := make(map[string]attr.Value)
		for lang, data := range statusPage.Translations {
			translationObject, diags := types.ObjectValue(
				translationSchema,
				map[string]attr.Value{
					"public_company_name": types.StringValue(data.PublicCompanyName),
					"header_logo_text":    types.StringValue(data.HeaderLogoText),
				},
			)
			translationData[lang] = translationObject
			diagnostics.Append(diags...)
			if diagnostics.HasError() {
				return nil
			}
		}
		// Create the translations map
		convertedTranslations, diags := types.MapValue(
			types.ObjectType{AttrTypes: translationSchema},
			translationData,
		)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return nil
		}

		translations = convertedTranslations
	}

	// When domain_config is active, the API mirrors its domain into the legacy domain field.
	// Keep domain empty to avoid conflicting with the planned value of "".
	legacyDomain := statusPage.Domain
	if statusPage.DomainConfig != nil {
		legacyDomain = ""
	}

	return &statusPageModel{
		Name:                           types.StringValue(statusPage.Name),
		Url:                            types.StringValue(statusPage.Url),
		TimeZone:                       types.StringValue(statusPage.TimeZone),
		Subdomain:                      types.StringValue(statusPage.Subdomain),
		SupportEmail:                   types.StringValue(statusPage.SupportEmail),
		TwitterPublicScreenName:        types.StringValue(statusPage.TwitterPublicScreenName),
		About:                          types.StringValue(statusPage.About),
		DisplayAbout:                   types.BoolValue(statusPage.DisplayAbout),
		CustomDomainEnabled:            types.BoolValue(statusPage.CustomDomainEnabled),
		Domain:                         types.StringValue(legacyDomain),
		DomainConfig:                   domainConfig,
		RestrictedIps:                  types.StringValue(statusPage.RestrictedIps),
		MemberRestricted:               types.BoolValue(statusPage.MemberRestricted),
		ScheduledMaintenanceDays:       types.Int64Value(statusPage.ScheduledMaintenanceDays),
		CustomJs:                       types.StringValue(statusPage.CustomJs),
		HeadCode:                       types.StringValue(statusPage.HeadCode),
		DateFormat:                     types.StringValue(statusPage.DateFormat),
		TimeFormat:                     types.StringValue(statusPage.TimeFormat),
		DateFormatEnforceEverywhere:    types.BoolValue(statusPage.DateFormatEnforceEverywhere),
		DisplayCalendar:                types.BoolValue(statusPage.DisplayCalendar),
		HideWatermark:                  types.BoolValue(statusPage.HideWatermark),
		MinorNotificationHours:         types.Int64Value(statusPage.MinorNotificationHours),
		MajorNotificationHours:         types.Int64Value(statusPage.MajorNotificationHours),
		MaintenanceNotificationHours:   types.Int64Value(statusPage.MaintenanceNotificationHours),
		HistoryLimitDays:               types.Int64Value(statusPage.HistoryLimitDays),
		CustomIncidentTypesEnabled:     types.BoolValue(statusPage.CustomIncidentTypesEnabled),
		InfoNoticesEnabled:             types.BoolValue(statusPage.InfoNoticesEnabled),
		LockedWhenMaintenance:          types.BoolValue(statusPage.LockedWhenMaintenance),
		Noindex:                        types.BoolValue(statusPage.Noindex),
		EnableAutoTranslations:         types.BoolValue(statusPage.EnableAutoTranslations),
		CaptchaEnabled:                 types.BoolValue(statusPage.CaptchaEnabled),
		Translations:                   translations,
		HeaderLogoText:                 types.StringValue(statusPage.HeaderLogoText),
		PublicCompanyName:              types.StringValue(statusPage.PublicCompanyName),
		BgImage:                        types.StringValue(statusPage.BgImage),
		Logo:                           types.StringValue(statusPage.Logo),
		Favicon:                        types.StringValue(statusPage.Favicon),
		DisplayUptimeGraph:             types.BoolValue(statusPage.DisplayUptimeGraph),
		UptimeGraphDays:                types.Int64Value(statusPage.UptimeGraphDays),
		CurrentIncidentsPosition:       types.StringValue(statusPage.CurrentIncidentsPosition),
		ThemeSelected:                  types.StringValue(statusPage.ThemeSelected),
		LinkColor:                      types.StringValue(statusPage.LinkColor),
		HeaderBgColor1:                 types.StringValue(statusPage.HeaderBgColor1),
		HeaderBgColor2:                 types.StringValue(statusPage.HeaderBgColor2),
		HeaderFgColor:                  types.StringValue(statusPage.HeaderFgColor),
		IncidentHeaderColor:            types.StringValue(statusPage.IncidentHeaderColor),
		IncidentLinkColor:              types.StringValue(statusPage.IncidentLinkColor),
		StatusOkColor:                  types.StringValue(statusPage.StatusOkColor),
		StatusMinorColor:               types.StringValue(statusPage.StatusMinorColor),
		StatusMajorColor:               types.StringValue(statusPage.StatusMajorColor),
		StatusMaintenanceColor:         types.StringValue(statusPage.StatusMaintenanceColor),
		CustomCss:                      types.StringValue(statusPage.CustomCss),
		CustomHeader:                   types.StringValue(statusPage.CustomHeader),
		CustomFooter:                   types.StringValue(statusPage.CustomFooter),
		NotifyByDefault:                types.BoolValue(statusPage.NotifyByDefault),
		TweetByDefault:                 types.BoolValue(statusPage.TweetByDefault),
		SlackSubscriptionsEnabled:      types.BoolValue(statusPage.SlackSubscriptionsEnabled),
		DiscordNotificationsEnabled:    types.BoolValue(statusPage.DiscordNotificationsEnabled),
		TeamsNotificationsEnabled:      types.BoolValue(statusPage.TeamsNotificationsEnabled),
		GoogleChatNotificationsEnabled: types.BoolValue(statusPage.GoogleChatNotificationsEnabled),
		MattermostNotificationsEnabled: types.BoolValue(statusPage.MattermostNotificationsEnabled),
		SmsNotificationsEnabled:        types.BoolValue(statusPage.SmsNotificationsEnabled),
		FeedEnabled:                    types.BoolValue(statusPage.FeedEnabled),
		CalendarEnabled:                types.BoolValue(statusPage.CalendarEnabled),
		GoogleCalendarEnabled:          types.BoolValue(statusPage.GoogleCalendarEnabled),
		SubscribersEnabled:             types.BoolValue(statusPage.SubscribersEnabled),
		NotificationEmail:              types.StringValue(statusPage.NotificationEmail),
		ReplyToEmail:                   types.StringValue(statusPage.ReplyToEmail),
		TweetingEnabled:                types.BoolValue(statusPage.TweetingEnabled),
		EmailLayoutTemplate:            types.StringValue(statusPage.EmailLayoutTemplate),
		EmailConfirmationTemplate:      types.StringValue(statusPage.EmailConfirmationTemplate),
		EmailNotificationTemplate:      types.StringValue(statusPage.EmailNotificationTemplate),
		EmailTemplatesEnabled:          types.BoolValue(statusPage.EmailTemplatesEnabled),
		ZoomNotificationsEnabled:       types.BoolValue(statusPage.ZoomNotificationsEnabled),
		AllowedEmailDomains:            types.StringValue(statusPage.AllowedEmailDomains),
		InsertedAt:                     types.StringValue(statusPage.InsertedAt),
		UpdatedAt:                      types.StringValue(statusPage.UpdatedAt),
	}
}

func stringPtrOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
