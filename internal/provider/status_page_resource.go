package provider

import (
	"context"
	"fmt"
	"strings"

	statuspal "terraform-provider-statuspal/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

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
	DisplayUptimeGraph             types.Bool   `tfsdk:"display_uptime_graph"`
	UptimeGraphDays                types.Int64  `tfsdk:"uptime_graph_days"`
	CurrentIncidentsPosition       types.String `tfsdk:"current_incidents_position"`
	ThemeSelected                  types.String `tfsdk:"theme_selected"`
	ThemeConfigs                   types.Object `tfsdk:"theme_configs"`
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
	InsertedAt                     types.String `tfsdk:"inserted_at"`
	UpdatedAt                      types.String `tfsdk:"updated_at"`
}

type statusPageThemeConfigsModel struct {
	LinkColor              types.String `tfsdk:"link_color"`
	HeaderBgColor1         types.String `tfsdk:"header_bg_color1"`
	HeaderBgColor2         types.String `tfsdk:"header_bg_color2"`
	HeaderFgColor          types.String `tfsdk:"header_fg_color"`
	IncidentHeaderColor    types.String `tfsdk:"incident_header_color"`
	StatusOkColor          types.String `tfsdk:"status_ok_color"`
	StatusMinorColor       types.String `tfsdk:"status_minor_color"`
	StatusMajorColor       types.String `tfsdk:"status_major_color"`
	StatusMaintenanceColor types.String `tfsdk:"status_maintenance_color"`
}

type statusPageTranslationsModel map[string]statusPageTranslationModel

type statusPageTranslationModel struct {
	PublicCompanyName types.String `tfsdk:"public_company_name"`
	HeaderLogoText    types.String `tfsdk:"header_logo_text"`
}

// Metadata returns the resource type name.
func (r *statusPageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
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
					},
					"twitter_public_screen_name": schema.StringAttribute{
						Description: "Twitter handle name (e.g. yourcompany).",
						Optional:    true,
						Computed:    true,
					},
					"about": schema.StringAttribute{
						Description: "Customize the about information displayed in your status page.",
						Optional:    true,
						Computed:    true,
					},
					"display_about": schema.BoolAttribute{
						Description: "Display about information.",
						Optional:    true,
						Computed:    true,
					},
					"custom_domain_enabled": schema.BoolAttribute{
						Description: "Enable your custom domain with SSL.",
						Optional:    true,
						Computed:    true,
					},
					"domain": schema.StringAttribute{
						Description: "Configure your own domain to point to your status page (e.g. status.your-company.com), we generate and auto-renew its SSL certificate for you.",
						Optional:    true,
						Computed:    true,
					},
					"restricted_ips": schema.StringAttribute{
						Description: `Your status page will be accessible only from this IPs (e.g. "1.1.1.1, 2.2.2.2").`,
						Optional:    true,
						Computed:    true,
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
					},
					"head_code": schema.StringAttribute{
						MarkdownDescription: "We'll insert this content inside the `<head>` tag.",
						Optional:            true,
						Computed:            true,
					},
					"date_format": schema.StringAttribute{
						Description: "Display timestamps of incidents and updates in this format.",
						Optional:    true,
						Computed:    true,
					},
					"time_format": schema.StringAttribute{
						Description: "Display timestamps of incidents and updates in this format.",
						Optional:    true,
						Computed:    true,
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
					},
					"bg_image": schema.StringAttribute{
						Description: "Background image url of the status page.",
						Optional:    true,
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
					"theme_configs": schema.ObjectAttribute{
						Description: "Theme configuration for the status page.",
						Optional:    true,
						Computed:    true,
						AttributeTypes: map[string]attr.Type{
							"link_color":               types.StringType,
							"header_bg_color1":         types.StringType,
							"header_bg_color2":         types.StringType,
							"header_fg_color":          types.StringType,
							"incident_header_color":    types.StringType,
							"status_ok_color":          types.StringType,
							"status_minor_color":       types.StringType,
							"status_major_color":       types.StringType,
							"status_maintenance_color": types.StringType,
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
					},
					"custom_header": schema.StringAttribute{
						MarkdownDescription: "A custom header for the status page (e.g. \"`<header>...</header>`\").",
						Optional:            true,
						Computed:            true,
					},
					"custom_footer": schema.StringAttribute{
						MarkdownDescription: "A custom footer for the status page (e.g. \"`<footer>...</footer>`\").",
						Optional:            true,
						Computed:            true,
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
					},
					"reply_to_email": schema.StringAttribute{
						Description: "The email address we'll use in the 'reply_to' field in emails to your subscribers. So they can reply to your notification emails.",
						Optional:    true,
						Computed:    true,
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
					},
					"email_confirmation_template": schema.StringAttribute{
						MarkdownDescription: "Custom confirmation email template, see the documentation: [Custom email templates](https://docs.statuspal.io/platform/subscriptions-and-notifications/custom-email-templates).",
						Optional:            true,
						Computed:            true,
					},
					"email_notification_template": schema.StringAttribute{
						MarkdownDescription: "Custom email notification template, see the documentation: [Custom email templates](https://docs.statuspal.io/platform/subscriptions-and-notifications/custom-email-templates).",
						Optional:            true,
						Computed:            true,
					},
					"email_templates_enabled": schema.BoolAttribute{
						Description: "The templates won't be used until this is enabled, but you can send test emails.",
						Optional:    true,
						Computed:    true,
					},
					"inserted_at": schema.StringAttribute{
						Description: "Datetime at which the status page was inserted.",
						Computed:    true,
					},
					"updated_at": schema.StringAttribute{
						Description: "Datetime at which the status page was last updated.",
						Computed:    true,
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

func (r *statusPageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
func (r *statusPageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*statuspal.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *statuspal.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func mapStatusPageModelToRequestBody(ctx *context.Context, statusPage *statusPageModel, diagnostics *diag.Diagnostics) *statuspal.StatusPage {
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

	var themeConfigs statusPageThemeConfigsModel
	if !statusPage.ThemeConfigs.IsNull() && !statusPage.ThemeConfigs.IsUnknown() {
		diags := statusPage.ThemeConfigs.As(*ctx, &themeConfigs, basetypes.ObjectAsOptions{})
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return nil
		}
	}

	return &statuspal.StatusPage{
		Name:                         statusPage.Name.ValueString(),
		Url:                          statusPage.Url.ValueString(),
		TimeZone:                     statusPage.TimeZone.ValueString(),
		Subdomain:                    statusPage.Subdomain.ValueString(),
		SupportEmail:                 statusPage.SupportEmail.ValueString(),
		TwitterPublicScreenName:      statusPage.TwitterPublicScreenName.ValueString(),
		About:                        statusPage.About.ValueString(),
		DisplayAbout:                 statusPage.DisplayAbout.ValueBool(),
		CustomDomainEnabled:          statusPage.CustomDomainEnabled.ValueBool(),
		Domain:                       statusPage.Domain.ValueString(),
		RestrictedIps:                statusPage.RestrictedIps.ValueString(),
		MemberRestricted:             statusPage.MemberRestricted.ValueBool(),
		ScheduledMaintenanceDays:     statusPage.ScheduledMaintenanceDays.ValueInt64(),
		CustomJs:                     statusPage.CustomJs.ValueString(),
		HeadCode:                     statusPage.HeadCode.ValueString(),
		DateFormat:                   statusPage.DateFormat.ValueString(),
		TimeFormat:                   statusPage.TimeFormat.ValueString(),
		DateFormatEnforceEverywhere:  statusPage.DateFormatEnforceEverywhere.ValueBool(),
		DisplayCalendar:              statusPage.DisplayCalendar.ValueBool(),
		HideWatermark:                statusPage.HideWatermark.ValueBool(),
		MinorNotificationHours:       statusPage.MinorNotificationHours.ValueInt64(),
		MajorNotificationHours:       statusPage.MajorNotificationHours.ValueInt64(),
		MaintenanceNotificationHours: statusPage.MaintenanceNotificationHours.ValueInt64(),
		HistoryLimitDays:             statusPage.HistoryLimitDays.ValueInt64(),
		CustomIncidentTypesEnabled:   statusPage.CustomIncidentTypesEnabled.ValueBool(),
		InfoNoticesEnabled:           statusPage.InfoNoticesEnabled.ValueBool(),
		LockedWhenMaintenance:        statusPage.LockedWhenMaintenance.ValueBool(),
		Noindex:                      statusPage.Noindex.ValueBool(),
		EnableAutoTranslations:       statusPage.EnableAutoTranslations.ValueBool(),
		CaptchaEnabled:               statusPage.CaptchaEnabled.ValueBool(),
		Translations:                 translationData,
		HeaderLogoText:               statusPage.HeaderLogoText.ValueString(),
		PublicCompanyName:            statusPage.PublicCompanyName.ValueString(),
		BgImage:                      statusPage.BgImage.ValueString(),
		DisplayUptimeGraph:           statusPage.DisplayUptimeGraph.ValueBool(),
		UptimeGraphDays:              statusPage.UptimeGraphDays.ValueInt64(),
		CurrentIncidentsPosition:     statusPage.CurrentIncidentsPosition.ValueString(),
		ThemeSelected:                statusPage.ThemeSelected.ValueString(),
		ThemeConfigs: statuspal.StatusPageThemeConfigs{
			LinkColor:              themeConfigs.LinkColor.ValueString(),
			HeaderBgColor1:         themeConfigs.HeaderBgColor1.ValueString(),
			HeaderBgColor2:         themeConfigs.HeaderBgColor2.ValueString(),
			HeaderFgColor:          themeConfigs.HeaderFgColor.ValueString(),
			IncidentHeaderColor:    themeConfigs.IncidentHeaderColor.ValueString(),
			StatusOkColor:          themeConfigs.StatusOkColor.ValueString(),
			StatusMinorColor:       themeConfigs.StatusMinorColor.ValueString(),
			StatusMajorColor:       themeConfigs.StatusMajorColor.ValueString(),
			StatusMaintenanceColor: themeConfigs.StatusMaintenanceColor.ValueString(),
		},
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
	}
}

func mapResponseToStatusPageModel(statusPage *statuspal.StatusPage, diagnostics *diag.Diagnostics) *statusPageModel {
	// Define the translation object schema
	translationSchema := map[string]attr.Type{
		"public_company_name": types.StringType,
		"header_logo_text":    types.StringType,
	}
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
	translations, diags := types.MapValue(
		types.ObjectType{AttrTypes: translationSchema},
		translationData,
	)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return nil
	}

	return &statusPageModel{
		Name:                         types.StringValue(statusPage.Name),
		Url:                          types.StringValue(statusPage.Url),
		TimeZone:                     types.StringValue(statusPage.TimeZone),
		Subdomain:                    types.StringValue(statusPage.Subdomain),
		SupportEmail:                 types.StringValue(statusPage.SupportEmail),
		TwitterPublicScreenName:      types.StringValue(statusPage.TwitterPublicScreenName),
		About:                        types.StringValue(statusPage.About),
		DisplayAbout:                 types.BoolValue(statusPage.DisplayAbout),
		CustomDomainEnabled:          types.BoolValue(statusPage.CustomDomainEnabled),
		Domain:                       types.StringValue(statusPage.Domain),
		RestrictedIps:                types.StringValue(statusPage.RestrictedIps),
		MemberRestricted:             types.BoolValue(statusPage.MemberRestricted),
		ScheduledMaintenanceDays:     types.Int64Value(statusPage.ScheduledMaintenanceDays),
		CustomJs:                     types.StringValue(statusPage.CustomJs),
		HeadCode:                     types.StringValue(statusPage.HeadCode),
		DateFormat:                   types.StringValue(statusPage.DateFormat),
		TimeFormat:                   types.StringValue(statusPage.TimeFormat),
		DateFormatEnforceEverywhere:  types.BoolValue(statusPage.DateFormatEnforceEverywhere),
		DisplayCalendar:              types.BoolValue(statusPage.DisplayCalendar),
		HideWatermark:                types.BoolValue(statusPage.HideWatermark),
		MinorNotificationHours:       types.Int64Value(statusPage.MinorNotificationHours),
		MajorNotificationHours:       types.Int64Value(statusPage.MajorNotificationHours),
		MaintenanceNotificationHours: types.Int64Value(statusPage.MaintenanceNotificationHours),
		HistoryLimitDays:             types.Int64Value(statusPage.HistoryLimitDays),
		CustomIncidentTypesEnabled:   types.BoolValue(statusPage.CustomIncidentTypesEnabled),
		InfoNoticesEnabled:           types.BoolValue(statusPage.InfoNoticesEnabled),
		LockedWhenMaintenance:        types.BoolValue(statusPage.LockedWhenMaintenance),
		Noindex:                      types.BoolValue(statusPage.Noindex),
		EnableAutoTranslations:       types.BoolValue(statusPage.EnableAutoTranslations),
		CaptchaEnabled:               types.BoolValue(statusPage.CaptchaEnabled),
		Translations:                 translations,
		HeaderLogoText:               types.StringValue(statusPage.HeaderLogoText),
		PublicCompanyName:            types.StringValue(statusPage.PublicCompanyName),
		BgImage:                      types.StringValue(statusPage.BgImage),
		DisplayUptimeGraph:           types.BoolValue(statusPage.DisplayUptimeGraph),
		UptimeGraphDays:              types.Int64Value(statusPage.UptimeGraphDays),
		CurrentIncidentsPosition:     types.StringValue(statusPage.CurrentIncidentsPosition),
		ThemeSelected:                types.StringValue(statusPage.ThemeSelected),
		ThemeConfigs: types.ObjectValueMust(map[string]attr.Type{
			"link_color":               types.StringType,
			"header_bg_color1":         types.StringType,
			"header_bg_color2":         types.StringType,
			"header_fg_color":          types.StringType,
			"incident_header_color":    types.StringType,
			"status_ok_color":          types.StringType,
			"status_minor_color":       types.StringType,
			"status_major_color":       types.StringType,
			"status_maintenance_color": types.StringType,
		}, map[string]attr.Value{
			"link_color":               types.StringValue(statusPage.ThemeConfigs.LinkColor),
			"header_bg_color1":         types.StringValue(statusPage.ThemeConfigs.HeaderBgColor1),
			"header_bg_color2":         types.StringValue(statusPage.ThemeConfigs.HeaderBgColor2),
			"header_fg_color":          types.StringValue(statusPage.ThemeConfigs.HeaderFgColor),
			"incident_header_color":    types.StringValue(statusPage.ThemeConfigs.IncidentHeaderColor),
			"status_ok_color":          types.StringValue(statusPage.ThemeConfigs.StatusOkColor),
			"status_minor_color":       types.StringValue(statusPage.ThemeConfigs.StatusMinorColor),
			"status_major_color":       types.StringValue(statusPage.ThemeConfigs.StatusMajorColor),
			"status_maintenance_color": types.StringValue(statusPage.ThemeConfigs.StatusMaintenanceColor),
		}),
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
		InsertedAt:                     types.StringValue(statusPage.InsertedAt),
		UpdatedAt:                      types.StringValue(statusPage.UpdatedAt),
	}
}
