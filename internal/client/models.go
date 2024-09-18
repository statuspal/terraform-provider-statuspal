package statuspal

// StatusPage struct.
type StatusPage struct {
	Name                           string                 `json:"name,omitempty"`
	Url                            string                 `json:"url,omitempty"`
	TimeZone                       string                 `json:"time_zone,omitempty"`
	Subdomain                      string                 `json:"subdomain,omitempty"`
	SupportEmail                   string                 `json:"support_email,omitempty"`
	TwitterPublicScreenName        string                 `json:"twitter_public_screen_name,omitempty"`
	About                          string                 `json:"about,omitempty"`
	DisplayAbout                   bool                   `json:"display_about"`
	CustomDomainEnabled            bool                   `json:"custom_domain_enabled"`
	Domain                         string                 `json:"domain,omitempty"`
	RestrictedIps                  string                 `json:"restricted_ips,omitempty"`
	MemberRestricted               bool                   `json:"member_restricted"`
	ScheduledMaintenanceDays       int64                  `json:"scheduled_maintenance_days,omitempty"`
	CustomJs                       string                 `json:"custom_js,omitempty"`
	HeadCode                       string                 `json:"head_code,omitempty"`
	DateFormat                     string                 `json:"date_format,omitempty"`
	TimeFormat                     string                 `json:"time_format,omitempty"`
	DateFormatEnforceEverywhere    bool                   `json:"date_format_enforce_everywhere"`
	DisplayCalendar                bool                   `json:"display_calendar"`
	HideWatermark                  bool                   `json:"hide_watermark"`
	MinorNotificationHours         int64                  `json:"minor_notification_hours,omitempty"`
	MajorNotificationHours         int64                  `json:"major_notification_hours,omitempty"`
	MaintenanceNotificationHours   int64                  `json:"maintenance_notification_hours,omitempty"`
	HistoryLimitDays               int64                  `json:"history_limit_days,omitempty"`
	CustomIncidentTypesEnabled     bool                   `json:"custom_incident_types_enabled"`
	InfoNoticesEnabled             bool                   `json:"info_notices_enabled"`
	LockedWhenMaintenance          bool                   `json:"locked_when_maintenance"`
	Noindex                        bool                   `json:"noindex"`
	EnableAutoTranslations         bool                   `json:"enable_auto_translations"`
	CaptchaEnabled                 bool                   `json:"captcha_enabled"`
	Translations                   StatusPageTranslations `json:"translations,omitempty"`
	HeaderLogoText                 string                 `json:"header_logo_text,omitempty"`
	PublicCompanyName              string                 `json:"public_company_name,omitempty"`
	BgImage                        string                 `json:"bg_image,omitempty"`
	Logo                           string                 `json:"logo,omitempty"`
	Favicon                        string                 `json:"favicon,omitempty"`
	DisplayUptimeGraph             bool                   `json:"display_uptime_graph"`
	UptimeGraphDays                int64                  `json:"uptime_graph_days,omitempty"`
	CurrentIncidentsPosition       string                 `json:"current_incidents_position,omitempty"`
	ThemeSelected                  string                 `json:"theme_selected,omitempty"`
	ThemeConfigs                   StatusPageThemeConfigs `json:"theme_configs,omitempty"`
	LinkColor                      string                 `json:"link_color,omitempty"`
	HeaderBgColor1                 string                 `json:"header_bg_color1,omitempty"`
	HeaderBgColor2                 string                 `json:"header_bg_color2,omitempty"`
	HeaderFgColor                  string                 `json:"header_fg_color,omitempty"`
	IncidentHeaderColor            string                 `json:"incident_header_color,omitempty"`
	IncidentLinkColor              string                 `json:"incident_link_color,omitempty"`
	StatusOkColor                  string                 `json:"status_ok_color,omitempty"`
	StatusMinorColor               string                 `json:"status_minor_color,omitempty"`
	StatusMajorColor               string                 `json:"status_major_color,omitempty"`
	StatusMaintenanceColor         string                 `json:"status_maintenance_color,omitempty"`
	CustomCss                      string                 `json:"custom_css,omitempty"`
	CustomHeader                   string                 `json:"custom_header,omitempty"`
	CustomFooter                   string                 `json:"custom_footer,omitempty"`
	NotifyByDefault                bool                   `json:"notify_by_default"`
	TweetByDefault                 bool                   `json:"tweet_by_default"`
	SlackSubscriptionsEnabled      bool                   `json:"slack_subscriptions_enabled"`
	DiscordNotificationsEnabled    bool                   `json:"discord_notifications_enabled"`
	TeamsNotificationsEnabled      bool                   `json:"teams_notifications_enabled"`
	GoogleChatNotificationsEnabled bool                   `json:"google_chat_notifications_enabled"`
	MattermostNotificationsEnabled bool                   `json:"mattermost_notifications_enabled"`
	SmsNotificationsEnabled        bool                   `json:"sms_notifications_enabled"`
	FeedEnabled                    bool                   `json:"feed_enabled"`
	CalendarEnabled                bool                   `json:"calendar_enabled"`
	GoogleCalendarEnabled          bool                   `json:"google_calendar_enabled"`
	SubscribersEnabled             bool                   `json:"subscribers_enabled"`
	NotificationEmail              string                 `json:"notification_email,omitempty"`
	ReplyToEmail                   string                 `json:"reply_to_email,omitempty"`
	TweetingEnabled                bool                   `json:"tweeting_enabled"`
	EmailLayoutTemplate            string                 `json:"email_layout_template,omitempty"`
	EmailConfirmationTemplate      string                 `json:"email_confirmation_template,omitempty"`
	EmailNotificationTemplate      string                 `json:"email_notification_template,omitempty"`
	EmailTemplatesEnabled          bool                   `json:"email_templates_enabled"`
	InsertedAt                     string                 `json:"inserted_at,omitempty"`
	UpdatedAt                      string                 `json:"updated_at,omitempty"`
}

type StatusPageTranslations map[string]StatusPageTranslation

type StatusPageTranslation struct {
	PublicCompanyName string `json:"public_company_name,omitempty"`
	HeaderLogoText    string `json:"header_logo_text,omitempty"`
}

type StatusPageThemeConfigs struct {
	LinkColor              string `json:"link_color,omitempty"`
	HeaderBgColor1         string `json:"header_bg_color1,omitempty"`
	HeaderBgColor2         string `json:"header_bg_color2,omitempty"`
	HeaderFgColor          string `json:"header_fg_color,omitempty"`
	IncidentHeaderColor    string `json:"incident_header_color,omitempty"`
	StatusOkColor          string `json:"status_ok_color,omitempty"`
	StatusMinorColor       string `json:"status_minor_color,omitempty"`
	StatusMajorColor       string `json:"status_major_color,omitempty"`
	StatusMaintenanceColor string `json:"status_maintenance_color,omitempty"`
}

// Service struct.
type Service struct {
	ID                                int64                         `json:"id,omitempty"`
	Name                              string                        `json:"name,omitempty"`
	Description                       string                        `json:"description,omitempty"`
	PrivateDescription                string                        `json:"private_description,omitempty"`
	ParentID                          int64                         `json:"parent_id,omitempty"`
	CurrentIncidentType               string                        `json:"current_incident_type,omitempty"`
	Monitoring                        string                        `json:"monitoring"`
	WebhookMonitoringService          string                        `json:"webhook_monitoring_service,omitempty"`
	WebhookCustomJsonpathSettings     WebhookCustomJsonpathSettings `json:"webhook_custom_jsonpath_settings,omitempty"`
	InboundEmailAddress               string                        `json:"inbound_email_address,omitempty"`
	IncomingWebhookUrl                string                        `json:"incoming_webhook_url,omitempty"`
	PingUrl                           string                        `json:"ping_url,omitempty"`
	IncidentType                      string                        `json:"incident_type,omitempty"`
	ParentIncidentType                string                        `json:"parent_incident_type,omitempty"`
	IsUp                              bool                          `json:"is_up"`
	PauseMonitoringDuringMaintenances bool                          `json:"pause_monitoring_during_maintenances"`
	InboundEmailID                    string                        `json:"inbound_email_id,omitempty"`
	AutoIncident                      bool                          `json:"auto_incident"`
	AutoNotify                        bool                          `json:"auto_notify"`
	ChildrenIDs                       []int64                       `json:"children_ids,omitempty"`
	Translations                      ServiceTranslations           `json:"translations,omitempty"`
	Private                           bool                          `json:"private"`
	DisplayUptimeGraph                bool                          `json:"display_uptime_graph"`
	DisplayResponseTimeChart          bool                          `json:"display_response_time_chart"`
	Order                             int64                         `json:"order,omitempty"`
	InsertedAt                        string                        `json:"inserted_at,omitempty"`
	UpdatedAt                         string                        `json:"updated_at,omitempty"`
}

type WebhookCustomJsonpathSettings struct {
	Jsonpath       string `json:"jsonpath,omitempty"`
	ExpectedResult string `json:"expected_result,omitempty"`
}

type ServiceTranslations map[string]ServiceTranslation

type ServiceTranslation struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}
