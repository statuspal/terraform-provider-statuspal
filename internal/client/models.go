package statuspal

// StatusPage struct.
type StatusPage struct {
	Name                           string                 `json:"name"`
	Url                            string                 `json:"url"`
	TimeZone                       string                 `json:"time_zone"`
	Subdomain                      string                 `json:"subdomain"`
	SupportEmail                   string                 `json:"support_email"`
	TwitterPublicScreenName        string                 `json:"twitter_public_screen_name"`
	About                          string                 `json:"about"`
	DisplayAbout                   bool                   `json:"display_about"`
	CustomDomainEnabled            bool                   `json:"custom_domain_enabled"`
	Domain                         string                 `json:"domain"`
	RestrictedIps                  string                 `json:"restricted_ips"`
	MemberRestricted               bool                   `json:"member_restricted"`
	ScheduledMaintenanceDays       int64                  `json:"scheduled_maintenance_days"`
	CustomJs                       string                 `json:"custom_js"`
	HeadCode                       string                 `json:"head_code"`
	DateFormat                     string                 `json:"date_format"`
	TimeFormat                     string                 `json:"time_format"`
	DateFormatEnforceEverywhere    bool                   `json:"date_format_enforce_everywhere"`
	DisplayCalendar                bool                   `json:"display_calendar"`
	HideWatermark                  bool                   `json:"hide_watermark"`
	MinorNotificationHours         int64                  `json:"minor_notification_hours"`
	MajorNotificationHours         int64                  `json:"major_notification_hours"`
	MaintenanceNotificationHours   int64                  `json:"maintenance_notification_hours"`
	HistoryLimitDays               int64                  `json:"history_limit_days"`
	CustomIncidentTypesEnabled     bool                   `json:"custom_incident_types_enabled"`
	InfoNoticesEnabled             bool                   `json:"info_notices_enabled"`
	LockedWhenMaintenance          bool                   `json:"locked_when_maintenance"`
	Noindex                        bool                   `json:"noindex"`
	EnableAutoTranslations         bool                   `json:"enable_auto_translations"`
	CaptchaEnabled                 bool                   `json:"captcha_enabled"`
	Translations                   StatusPageTranslations `json:"translations"`
	HeaderLogoText                 string                 `json:"header_logo_text"`
	PublicCompanyName              string                 `json:"public_company_name"`
	BgImage                        string                 `json:"bg_image"`
	Logo                           string                 `json:"logo"`
	Favicon                        string                 `json:"favicon"`
	DisplayUptimeGraph             bool                   `json:"display_uptime_graph"`
	UptimeGraphDays                int64                  `json:"uptime_graph_days"`
	CurrentIncidentsPosition       string                 `json:"current_incidents_position"`
	ThemeSelected                  string                 `json:"theme_selected"`
	LinkColor                      string                 `json:"link_color"`
	HeaderBgColor1                 string                 `json:"header_bg_color1"`
	HeaderBgColor2                 string                 `json:"header_bg_color2"`
	HeaderFgColor                  string                 `json:"header_fg_color"`
	IncidentHeaderColor            string                 `json:"incident_header_color"`
	IncidentLinkColor              string                 `json:"incident_link_color"`
	StatusOkColor                  string                 `json:"status_ok_color"`
	StatusMinorColor               string                 `json:"status_minor_color"`
	StatusMajorColor               string                 `json:"status_major_color"`
	StatusMaintenanceColor         string                 `json:"status_maintenance_color"`
	CustomCss                      string                 `json:"custom_css"`
	CustomHeader                   string                 `json:"custom_header"`
	CustomFooter                   string                 `json:"custom_footer"`
	NotifyByDefault                bool                   `json:"notify_by_default"`
	TweetByDefault                 bool                   `json:"tweet_by_default"`
	SlackSubscriptionsEnabled      bool                   `json:"slack_subscriptions_enabled"`
	DiscordNotificationsEnabled    bool                   `json:"discord_notifications_enabled"`
	TeamsNotificationsEnabled      bool                   `json:"teams_notifications_enabled"`
	ZoomNotificationsEnabled       bool                   `json:"zoom_notifications_enabled"`
	GoogleChatNotificationsEnabled bool                   `json:"google_chat_notifications_enabled"`
	MattermostNotificationsEnabled bool                   `json:"mattermost_notifications_enabled"`
	SmsNotificationsEnabled        bool                   `json:"sms_notifications_enabled"`
	AllowedEmailDomains            string                 `json:"allowed_email_domains"`
	FeedEnabled                    bool                   `json:"feed_enabled"`
	CalendarEnabled                bool                   `json:"calendar_enabled"`
	GoogleCalendarEnabled          bool                   `json:"google_calendar_enabled"`
	SubscribersEnabled             bool                   `json:"subscribers_enabled"`
	NotificationEmail              string                 `json:"notification_email"`
	ReplyToEmail                   string                 `json:"reply_to_email"`
	TweetingEnabled                bool                   `json:"tweeting_enabled"`
	EmailLayoutTemplate            string                 `json:"email_layout_template"`
	EmailConfirmationTemplate      string                 `json:"email_confirmation_template"`
	EmailNotificationTemplate      string                 `json:"email_notification_template"`
	EmailTemplatesEnabled          bool                   `json:"email_templates_enabled"`
	InsertedAt                     string                 `json:"inserted_at"`
	UpdatedAt                      string                 `json:"updated_at"`
}

type StatusPageTranslations map[string]StatusPageTranslation

type StatusPageTranslation struct {
	PublicCompanyName string `json:"public_company_name"`
	HeaderLogoText    string `json:"header_logo_text"`
}

// Service struct.
type Service struct {
	ID                                int64                          `json:"id"`
	Name                              string                         `json:"name"`
	Description                       string                         `json:"description"`
	PrivateDescription                string                         `json:"private_description"`
	ParentID                          *int64                         `json:"parent_id"`
	CurrentIncidentType               string                         `json:"current_incident_type"`
	Monitoring                        string                         `json:"monitoring"`
	WebhookMonitoringService          string                         `json:"webhook_monitoring_service"`
	WebhookCustomJsonpathSettings     *WebhookCustomJsonpathSettings `json:"webhook_custom_jsonpath_settings"`
	InboundEmailAddress               string                         `json:"inbound_email_address"`
	IncomingWebhookUrl                string                         `json:"incoming_webhook_url"`
	PingUrl                           string                         `json:"ping_url"`
	IncidentType                      string                         `json:"incident_type"`
	ParentIncidentType                string                         `json:"parent_incident_type"`
	IsUp                              bool                           `json:"is_up"`
	PauseMonitoringDuringMaintenances bool                           `json:"pause_monitoring_during_maintenances"`
	InboundEmailID                    string                         `json:"inbound_email_id"`
	AutoIncident                      bool                           `json:"auto_incident"`
	AutoNotify                        bool                           `json:"auto_notify"`
	ChildrenIDs                       []int64                        `json:"children_ids"`
	Translations                      ServiceTranslations            `json:"translations"`
	Private                           bool                           `json:"private"`
	DisplayUptimeGraph                bool                           `json:"display_uptime_graph"`
	DisplayResponseTimeChart          bool                           `json:"display_response_time_chart"`
	Order                             int64                          `json:"order"`
	MonitoringOptions                 *MonitoringOptions             `json:"monitoring_options"`
	InsertedAt                        string                         `json:"inserted_at"`
	UpdatedAt                         string                         `json:"updated_at"`
}

type WebhookCustomJsonpathSettings struct {
	Jsonpath       string `json:"jsonpath"`
	ExpectedResult string `json:"expected_result"`
}

type ServiceTranslations map[string]ServiceTranslation

type ServiceTranslation struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// MonitoringOptions struct defines the monitoring configuration for services.
type MonitoringOptions struct {
	Method                  string                   `json:"method"`
	Headers                 MonitoringOptionsHeaders `json:"headers"`
	KeywordDown             string                   `json:"keyword_down"`
	KeywordUp               string                   `json:"keyword_up"`
	ExternalServiceStatuses []string                 `json:"external_service_statuses"`
}

type MonitoringOptionsHeaders []MonitoringOptionsHeader
type MonitoringOptionsHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type NotificationRecipient struct {
	ID    int64  `json:"id"`
	Email string `json:"email"` // Add more fields as needed
}

// Metric represents a metric on the status page.
type Metric struct {
	ID              int64  `json:"id"`
	Status          string `json:"status"`
	LatestEntryTime int64  `json:"latest_entry_time"`
	Order           int64  `json:"order"`
	Title           string `json:"title"`
	Unit            string `json:"unit"`
	Type            string `json:"type"`
	Enabled         bool   `json:"enabled"`
	Visible         bool   `json:"visible"`
	RemoteID        string `json:"remote_id"`
	RemoteName      string `json:"remote_name"`
	Threshold       int64  `json:"threshold"`
	FeaturedNumber  string `json:"featured_number"`
	IntegrationID   *int64 `json:"integration_id"`
}
