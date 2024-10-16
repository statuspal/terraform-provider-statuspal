package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	statuspal "terraform-provider-statuspal/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &statusPagesDataSource{}
	_ datasource.DataSourceWithConfigure = &statusPagesDataSource{}
)

// NewStatusPagesDataSource is a helper function to simplify the provider implementation.
func NewStatusPagesDataSource() datasource.DataSource {
	return &statusPagesDataSource{}
}

// statusPagesDataSource is the data source implementation.
type statusPagesDataSource struct {
	client *statuspal.Client
}

// statusPagesDataSourceModel maps the data source schema data.
type statusPagesDataSourceModel struct {
	ID             types.String       `tfsdk:"id"` // only for test case
	OrganizationID types.String       `tfsdk:"organization_id"`
	StatusPages    []statusPagesModel `tfsdk:"status_pages"`
}

// statusPagesModel maps status_pages schema data.
type statusPagesModel struct {
	Name                           types.String                 `tfsdk:"name"`
	Url                            types.String                 `tfsdk:"url"`
	TimeZone                       types.String                 `tfsdk:"time_zone"`
	Subdomain                      types.String                 `tfsdk:"subdomain"`
	SupportEmail                   types.String                 `tfsdk:"support_email"`
	TwitterPublicScreenName        types.String                 `tfsdk:"twitter_public_screen_name"`
	About                          types.String                 `tfsdk:"about"`
	DisplayAbout                   types.Bool                   `tfsdk:"display_about"`
	CustomDomainEnabled            types.Bool                   `tfsdk:"custom_domain_enabled"`
	Domain                         types.String                 `tfsdk:"domain"`
	RestrictedIps                  types.String                 `tfsdk:"restricted_ips"`
	MemberRestricted               types.Bool                   `tfsdk:"member_restricted"`
	ScheduledMaintenanceDays       types.Int64                  `tfsdk:"scheduled_maintenance_days"`
	CustomJs                       types.String                 `tfsdk:"custom_js"`
	HeadCode                       types.String                 `tfsdk:"head_code"`
	DateFormat                     types.String                 `tfsdk:"date_format"`
	TimeFormat                     types.String                 `tfsdk:"time_format"`
	DateFormatEnforceEverywhere    types.Bool                   `tfsdk:"date_format_enforce_everywhere"`
	DisplayCalendar                types.Bool                   `tfsdk:"display_calendar"`
	HideWatermark                  types.Bool                   `tfsdk:"hide_watermark"`
	MinorNotificationHours         types.Int64                  `tfsdk:"minor_notification_hours"`
	MajorNotificationHours         types.Int64                  `tfsdk:"major_notification_hours"`
	MaintenanceNotificationHours   types.Int64                  `tfsdk:"maintenance_notification_hours"`
	HistoryLimitDays               types.Int64                  `tfsdk:"history_limit_days"`
	CustomIncidentTypesEnabled     types.Bool                   `tfsdk:"custom_incident_types_enabled"`
	InfoNoticesEnabled             types.Bool                   `tfsdk:"info_notices_enabled"`
	LockedWhenMaintenance          types.Bool                   `tfsdk:"locked_when_maintenance"`
	Noindex                        types.Bool                   `tfsdk:"noindex"`
	EnableAutoTranslations         types.Bool                   `tfsdk:"enable_auto_translations"`
	CaptchaEnabled                 types.Bool                   `tfsdk:"captcha_enabled"`
	Translations                   statusPagesTranslationsModel `tfsdk:"translations"`
	HeaderLogoText                 types.String                 `tfsdk:"header_logo_text"`
	PublicCompanyName              types.String                 `tfsdk:"public_company_name"`
	BgImage                        types.String                 `tfsdk:"bg_image"`
	Logo                           types.String                 `tfsdk:"logo"`
	Favicon                        types.String                 `tfsdk:"favicon"`
	DisplayUptimeGraph             types.Bool                   `tfsdk:"display_uptime_graph"`
	UptimeGraphDays                types.Int64                  `tfsdk:"uptime_graph_days"`
	CurrentIncidentsPosition       types.String                 `tfsdk:"current_incidents_position"`
	ThemeSelected                  types.String                 `tfsdk:"theme_selected"`
	ThemeConfigs                   statusPagesThemeConfigsModel `tfsdk:"theme_configs"`
	LinkColor                      types.String                 `tfsdk:"link_color"`
	HeaderBgColor1                 types.String                 `tfsdk:"header_bg_color1"`
	HeaderBgColor2                 types.String                 `tfsdk:"header_bg_color2"`
	HeaderFgColor                  types.String                 `tfsdk:"header_fg_color"`
	IncidentHeaderColor            types.String                 `tfsdk:"incident_header_color"`
	IncidentLinkColor              types.String                 `tfsdk:"incident_link_color"`
	StatusOkColor                  types.String                 `tfsdk:"status_ok_color"`
	StatusMinorColor               types.String                 `tfsdk:"status_minor_color"`
	StatusMajorColor               types.String                 `tfsdk:"status_major_color"`
	StatusMaintenanceColor         types.String                 `tfsdk:"status_maintenance_color"`
	CustomCss                      types.String                 `tfsdk:"custom_css"`
	CustomHeader                   types.String                 `tfsdk:"custom_header"`
	CustomFooter                   types.String                 `tfsdk:"custom_footer"`
	NotifyByDefault                types.Bool                   `tfsdk:"notify_by_default"`
	TweetByDefault                 types.Bool                   `tfsdk:"tweet_by_default"`
	SlackSubscriptionsEnabled      types.Bool                   `tfsdk:"slack_subscriptions_enabled"`
	DiscordNotificationsEnabled    types.Bool                   `tfsdk:"discord_notifications_enabled"`
	TeamsNotificationsEnabled      types.Bool                   `tfsdk:"teams_notifications_enabled"`
	GoogleChatNotificationsEnabled types.Bool                   `tfsdk:"google_chat_notifications_enabled"`
	MattermostNotificationsEnabled types.Bool                   `tfsdk:"mattermost_notifications_enabled"`
	SmsNotificationsEnabled        types.Bool                   `tfsdk:"sms_notifications_enabled"`
	FeedEnabled                    types.Bool                   `tfsdk:"feed_enabled"`
	CalendarEnabled                types.Bool                   `tfsdk:"calendar_enabled"`
	GoogleCalendarEnabled          types.Bool                   `tfsdk:"google_calendar_enabled"`
	SubscribersEnabled             types.Bool                   `tfsdk:"subscribers_enabled"`
	NotificationEmail              types.String                 `tfsdk:"notification_email"`
	ReplyToEmail                   types.String                 `tfsdk:"reply_to_email"`
	TweetingEnabled                types.Bool                   `tfsdk:"tweeting_enabled"`
	EmailLayoutTemplate            types.String                 `tfsdk:"email_layout_template"`
	EmailConfirmationTemplate      types.String                 `tfsdk:"email_confirmation_template"`
	EmailNotificationTemplate      types.String                 `tfsdk:"email_notification_template"`
	EmailTemplatesEnabled          types.Bool                   `tfsdk:"email_templates_enabled"`
	ZoomNotificationsEnabled       types.Bool                   `tfsdk:"zoom_notifications_enabled"`
	AllowedEmailDomains            types.String                 `tfsdk:"allowed_email_domains"`
	InsertedAt                     types.String                 `tfsdk:"inserted_at"`
	UpdatedAt                      types.String                 `tfsdk:"updated_at"`
}

type statusPagesThemeConfigsModel struct {
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

type statusPagesTranslationsModel map[string]statusPagesTranslationModel

type statusPagesTranslationModel struct {
	PublicCompanyName types.String `tfsdk:"public_company_name"`
	HeaderLogoText    types.String `tfsdk:"header_logo_text"`
}

// Metadata returns the data source type name.
func (d *statusPagesDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_status_pages"
}

// Schema defines the schema for the data source.
func (d *statusPagesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of status pages in the organization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier attribute. Ignore it, only used in testing.",
				Computed:    true,
			},
			"organization_id": schema.StringAttribute{
				Description: "The organization ID of the status pages.",
				Required:    true,
			},
			"status_pages": schema.ListNestedAttribute{
				Description: "List of status pages.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Company, project or service name.",
							Computed:    true,
						},
						"url": schema.StringAttribute{
							Description: "The website to your company, project or service.",
							Computed:    true,
						},
						"time_zone": schema.StringAttribute{
							Description: `The primary timezone the status page uses to display incidents (e.g. "Europe/Berlin").`,
							Computed:    true,
						},
						"subdomain": schema.StringAttribute{
							Description: "The status page subdomain on statuspal.",
							Computed:    true,
						},
						"support_email": schema.StringAttribute{
							Description: "Your company's support email.",
							Computed:    true,
						},
						"twitter_public_screen_name": schema.StringAttribute{
							Description: "Twitter handle name (e.g. yourcompany).",
							Computed:    true,
						},
						"about": schema.StringAttribute{
							Description: "Customize the about information displayed in your status page.",
							Computed:    true,
						},
						"display_about": schema.BoolAttribute{
							Description: "Display about information.",
							Computed:    true,
						},
						"custom_domain_enabled": schema.BoolAttribute{
							Description: "Enable your custom domain with SSL.",
							Computed:    true,
						},
						"domain": schema.StringAttribute{
							Description: "Configure your own domain to point to your status page (e.g. status.your-company.com), we generate and auto-renew its SSL certificate for you.",
							Computed:    true,
						},
						"restricted_ips": schema.StringAttribute{
							Description: `Your status page will be accessible only from this IPs (e.g. "1.1.1.1, 2.2.2.2").`,
							Computed:    true,
						},
						"member_restricted": schema.BoolAttribute{
							Description: "Only signed in members will be allowed to access your status page.",
							Computed:    true,
						},
						"scheduled_maintenance_days": schema.Int64Attribute{
							Description: "Display scheduled maintenance.",
							Computed:    true,
						},
						"custom_js": schema.StringAttribute{
							MarkdownDescription: "We'll insert this content inside the `<script>` tag at the bottom of your status page `<body>` tag.",
							Computed:            true,
						},
						"head_code": schema.StringAttribute{
							MarkdownDescription: "We'll insert this content inside the `<head>` tag.",
							Computed:            true,
						},
						"date_format": schema.StringAttribute{
							Description: "Display timestamps of incidents and updates in this format.",
							Computed:    true,
						},
						"time_format": schema.StringAttribute{
							Description: "Display timestamps of incidents and updates in this format.",
							Computed:    true,
						},
						"date_format_enforce_everywhere": schema.BoolAttribute{
							Description: "The above date format will be used everywhere in the status page. Timezone conversion to client's will be disabled.",
							Computed:    true,
						},
						"display_calendar": schema.BoolAttribute{
							Description: "Display uptime calendar at status page.",
							Computed:    true,
						},
						"hide_watermark": schema.BoolAttribute{
							Description: `Hide "Powered by Statuspal.io".`,
							Computed:    true,
						},
						"minor_notification_hours": schema.Int64Attribute{
							Description: "Long-running incident notification (Minor incident).",
							Computed:    true,
						},
						"major_notification_hours": schema.Int64Attribute{
							Description: "Long-running incident notification (Major incident).",
							Computed:    true,
						},
						"maintenance_notification_hours": schema.Int64Attribute{
							Description: "Long-running incident notification (Maintenance).",
							Computed:    true,
						},
						"history_limit_days": schema.Int64Attribute{
							Description: "Incident history limit (omit for No Limit).",
							Computed:    true,
						},
						"custom_incident_types_enabled": schema.BoolAttribute{
							Description: "Enable custom incident types.",
							Computed:    true,
						},
						"info_notices_enabled": schema.BoolAttribute{
							Description: "Enable information notices.",
							Computed:    true,
						},
						"locked_when_maintenance": schema.BoolAttribute{
							Description: "Lock from adding incidents when under maintenance.",
							Computed:    true,
						},
						"noindex": schema.BoolAttribute{
							Description: "Remove status page from being indexed by search engines (e.g. Google).",
							Computed:    true,
						},
						"enable_auto_translations": schema.BoolAttribute{
							Description: "Enable auto translations when creating incidents, maintenances and info notices.",
							Computed:    true,
						},
						"captcha_enabled": schema.BoolAttribute{
							Description: "Enable captchas (this option is only available when the status page is member restricted).",
							Computed:    true,
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
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"public_company_name": schema.StringAttribute{
										Description: "Displayed at the footer of the status page.",
										Computed:    true,
									},
									"header_logo_text": schema.StringAttribute{
										Description: "Displayed at the header of the status page.",
										Computed:    true,
									},
								},
							},
						},
						"header_logo_text": schema.StringAttribute{
							Description: "Displayed at the header of the status page.",
							Computed:    true,
						},
						"public_company_name": schema.StringAttribute{
							Description: "Displayed at the footer of the status page.",
							Computed:    true,
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
							Computed:    true,
						},
						"uptime_graph_days": schema.Int64Attribute{
							Description: "Uptime graph period.",
							Computed:    true,
						},
						"current_incidents_position": schema.StringAttribute{
							Description: `The incident position displayed in the status page, it can be "below_services" and "above_services".`,
							Computed:    true,
						},
						"theme_selected": schema.StringAttribute{
							Description: `The selected theme for state page, it can be "default" and "big-logo".`,
							Computed:    true,
						},
						"theme_configs": schema.ObjectAttribute{
							Description: "Theme configuration for the status page.",
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
							Computed:    true,
						},
						"header_bg_color1": schema.StringAttribute{
							Description: "The background color at left side of the status page header.",
							Computed:    true,
						},
						"header_bg_color2": schema.StringAttribute{
							Description: "The background color at right side of the status page header.",
							Computed:    true,
						},
						"header_fg_color": schema.StringAttribute{
							Description: "The text color in the status page.",
							Computed:    true,
						},
						"incident_header_color": schema.StringAttribute{
							Description: "Incidents header color in the status page.",
							Computed:    true,
						},
						"incident_link_color": schema.StringAttribute{
							Description: "Incidents link color in the status page.",
							Computed:    true,
						},
						"status_ok_color": schema.StringAttribute{
							Description: "The status page colors when there is no incident.",
							Computed:    true,
						},
						"status_minor_color": schema.StringAttribute{
							Description: "The status page colors when there is a minor incident.",
							Computed:    true,
						},
						"status_major_color": schema.StringAttribute{
							Description: "The status page colors when there is a major incident.",
							Computed:    true,
						},
						"status_maintenance_color": schema.StringAttribute{
							Description: "The status page colors when there is a maintenance incident.",
							Computed:    true,
						},
						"custom_css": schema.StringAttribute{
							MarkdownDescription: "We'll insert this content inside the `<style>` tag.",
							Computed:            true,
						},
						"custom_header": schema.StringAttribute{
							MarkdownDescription: "A custom header for the status page (e.g. \"`<header>...</header>`\").",
							Computed:            true,
						},
						"custom_footer": schema.StringAttribute{
							MarkdownDescription: "A custom footer for the status page (e.g. \"`<footer>...</footer>`\").",
							Computed:            true,
						},
						"notify_by_default": schema.BoolAttribute{
							Description: "Check the Notify subscribers checkbox by default.",
							Computed:    true,
						},
						"tweet_by_default": schema.BoolAttribute{
							Description: "Check the Tweet checkbox by default.",
							Computed:    true,
						},
						"slack_subscriptions_enabled": schema.BoolAttribute{
							Description: "Allow your customers to subscribe via Slack to updates on your status page's status.",
							Computed:    true,
						},
						"discord_notifications_enabled": schema.BoolAttribute{
							Description: "Allow your customers to receive notifications on a Discord channel.",
							Computed:    true,
						},
						"teams_notifications_enabled": schema.BoolAttribute{
							Description: "Allow your customers to receive notifications on Microsoft Teams.",
							Computed:    true,
						},
						"google_chat_notifications_enabled": schema.BoolAttribute{
							Description: "Allow your customers to receive notifications on Google Chat.",
							Computed:    true,
						},
						"mattermost_notifications_enabled": schema.BoolAttribute{
							Description: "Allow your customers to receive notifications on Mattermost.",
							Computed:    true,
						},
						"sms_notifications_enabled": schema.BoolAttribute{
							Description: "Allow your customers to receive SMS notifications on your status page's status (to enable this you need to have a Twilio or Esendex integration).",
							Computed:    true,
						},
						"feed_enabled": schema.BoolAttribute{
							Description: "Allow your customers to receive updates as RSS and Atom feeds.",
							Computed:    true,
						},
						"calendar_enabled": schema.BoolAttribute{
							Description: "Allow your customers to receive updates via iCalendar feed.",
							Computed:    true,
						},
						"google_calendar_enabled": schema.BoolAttribute{
							Description: "Allow your customers to import Google Calendar with Status Pages maintenance (business only).",
							Computed:    true,
						},
						"subscribers_enabled": schema.BoolAttribute{
							Description: "Allow email customers to receive email notifications.",
							Computed:    true,
						},
						"notification_email": schema.StringAttribute{
							Description: "Allow your customers to subscribe via email to updates on your status page's status.",
							Computed:    true,
						},
						"reply_to_email": schema.StringAttribute{
							Description: "The email address we'll use in the 'reply_to' field in emails to your subscribers. So they can reply to your notification emails.",
							Computed:    true,
						},
						"tweeting_enabled": schema.BoolAttribute{
							Description: "Allows to send tweets when creating or updating an incident.",
							Computed:    true,
						},
						"email_layout_template": schema.StringAttribute{
							MarkdownDescription: "Custom email layout template, see the documentation: [Custom email templates](https://docs.statuspal.io/platform/subscriptions-and-notifications/custom-email-templates).",
							Computed:            true,
						},
						"email_confirmation_template": schema.StringAttribute{
							MarkdownDescription: "Custom confirmation email template, see the documentation: [Custom email templates](https://docs.statuspal.io/platform/subscriptions-and-notifications/custom-email-templates).",
							Computed:            true,
						},
						"email_notification_template": schema.StringAttribute{
							MarkdownDescription: "Custom email notification template, see the documentation: [Custom email templates](https://docs.statuspal.io/platform/subscriptions-and-notifications/custom-email-templates).",
							Computed:            true,
						},
						"email_templates_enabled": schema.BoolAttribute{
							Description: "The templates won't be used until this is enabled, but you can send test emails.",
							Computed:    true,
						},
						"zoom_notifications_enabled": schema.BoolAttribute{
							Description: "Allow your customers to receive notifications on Zoom.",
							Computed:    true,
						},
						"allowed_email_domains": schema.StringAttribute{
							Description: "Allowed email domains. Each domain should be separated by \\n (e.g., 'acme.corp\\nnapster.com')",
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
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *statusPagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Retrieve values from config
	var state statusPagesDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationID := state.OrganizationID.ValueString()
	statusPages, err := d.client.GetStatusPages(&organizationID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read StatusPal StatusPages",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, statusPage := range *statusPages {
		// Create the translationData object dynamically
		translationData := make(statusPagesTranslationsModel)
		for lang, data := range statusPage.Translations {
			translationData[lang] = statusPagesTranslationModel{
				PublicCompanyName: types.StringValue(data.PublicCompanyName),
				HeaderLogoText:    types.StringValue(data.HeaderLogoText),
			}
		}

		statusPageState := statusPagesModel{
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
			Translations:                 translationData,
			HeaderLogoText:               types.StringValue(statusPage.HeaderLogoText),
			PublicCompanyName:            types.StringValue(statusPage.PublicCompanyName),
			BgImage:                      types.StringValue(statusPage.BgImage),
			Logo:                         types.StringValue(statusPage.Logo),
			Favicon:                      types.StringValue(statusPage.Favicon),
			DisplayUptimeGraph:           types.BoolValue(statusPage.DisplayUptimeGraph),
			UptimeGraphDays:              types.Int64Value(statusPage.UptimeGraphDays),
			CurrentIncidentsPosition:     types.StringValue(statusPage.CurrentIncidentsPosition),
			ThemeSelected:                types.StringValue(statusPage.ThemeSelected),
			ThemeConfigs: statusPagesThemeConfigsModel{
				LinkColor:              types.StringValue(statusPage.ThemeConfigs.LinkColor),
				HeaderBgColor1:         types.StringValue(statusPage.ThemeConfigs.HeaderBgColor1),
				HeaderBgColor2:         types.StringValue(statusPage.ThemeConfigs.HeaderBgColor2),
				HeaderFgColor:          types.StringValue(statusPage.ThemeConfigs.HeaderFgColor),
				IncidentHeaderColor:    types.StringValue(statusPage.ThemeConfigs.IncidentHeaderColor),
				StatusOkColor:          types.StringValue(statusPage.ThemeConfigs.StatusOkColor),
				StatusMinorColor:       types.StringValue(statusPage.ThemeConfigs.StatusMinorColor),
				StatusMajorColor:       types.StringValue(statusPage.ThemeConfigs.StatusMajorColor),
				StatusMaintenanceColor: types.StringValue(statusPage.ThemeConfigs.StatusMaintenanceColor),
			},
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

		state.StatusPages = append(state.StatusPages, statusPageState)
	}
	state.ID = types.StringValue("placeholder") // only for test case

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *statusPagesDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
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

	d.client = client
}
