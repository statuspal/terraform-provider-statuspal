package provider

import (
	"context"
	"fmt"
	"strconv"

	statuspal "terraform-provider-statuspal/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &servicesDataSource{}
	_ datasource.DataSourceWithConfigure = &servicesDataSource{}
)

// NewServicesDataSource is a helper function to simplify the provider implementation.
func NewServicesDataSource() datasource.DataSource {
	return &servicesDataSource{}
}

// servicesDataSource is the data source implementation.
type servicesDataSource struct {
	client *statuspal.Client
}

// servicesDataSourceModel maps the data source schema data.
type servicesDataSourceModel struct {
	ID                  types.String    `tfsdk:"id"` // only for test case
	StatusPageSubdomain types.String    `tfsdk:"status_page_subdomain"`
	Services            []servicesModel `tfsdk:"services"`
}

// servicesModel maps services schema data.
type servicesModel struct {
	ID                                types.String                               `tfsdk:"id"`
	Name                              types.String                               `tfsdk:"name"`
	Description                       types.String                               `tfsdk:"description"`
	PrivateDescription                types.String                               `tfsdk:"private_description"`
	ParentID                          types.String                               `tfsdk:"parent_id"`
	CurrentIncidentType               types.String                               `tfsdk:"current_incident_type"`
	Monitoring                        types.String                               `tfsdk:"monitoring"`
	WebhookMonitoringService          types.String                               `tfsdk:"webhook_monitoring_service"`
	WebhookCustomJsonpathSettings     servicesWebhookCustomJsonpathSettingsModel `tfsdk:"webhook_custom_jsonpath_settings"`
	InboundEmailAddress               types.String                               `tfsdk:"inbound_email_address"`
	IncomingWebhookUrl                types.String                               `tfsdk:"incoming_webhook_url"`
	PingUrl                           types.String                               `tfsdk:"ping_url"`
	IncidentType                      types.String                               `tfsdk:"incident_type"`
	ParentIncidentType                types.String                               `tfsdk:"parent_incident_type"`
	IsUp                              types.Bool                                 `tfsdk:"is_up"`
	PauseMonitoringDuringMaintenances types.Bool                                 `tfsdk:"pause_monitoring_during_maintenances"`
	InboundEmailID                    types.String                               `tfsdk:"inbound_email_id"`
	AutoIncident                      types.Bool                                 `tfsdk:"auto_incident"`
	AutoNotify                        types.Bool                                 `tfsdk:"auto_notify"`
	ChildrenIDs                       types.List                                 `tfsdk:"children_ids"`
	Translations                      servicesTranslationsModel                  `tfsdk:"translations"`
	Private                           types.Bool                                 `tfsdk:"private"`
	DisplayUptimeGraph                types.Bool                                 `tfsdk:"display_uptime_graph"`
	DisplayResponseTimeChart          types.Bool                                 `tfsdk:"display_response_time_chart"`
	Order                             types.Int64                                `tfsdk:"order"`
	InsertedAt                        types.String                               `tfsdk:"inserted_at"`
	UpdatedAt                         types.String                               `tfsdk:"updated_at"`
}

type servicesWebhookCustomJsonpathSettingsModel struct {
	Jsonpath       types.String `tfsdk:"jsonpath"`
	ExpectedResult types.String `tfsdk:"expected_result"`
}

type servicesTranslationsModel map[string]servicesTranslationModel

type servicesTranslationModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

// Metadata returns the data source type name.
func (d *servicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_services"
}

// Schema defines the schema for the data source.
func (d *servicesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of services in the status page.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier attribute. Ignore it, only used in testing.",
				Computed:    true,
			},
			"status_page_subdomain": schema.StringAttribute{
				Description: "The status page subdomain of the services.",
				Required:    true,
			},
			"services": schema.ListNestedAttribute{
				Description: "List of services.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The ID of the service.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the service.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "The description of the service.",
							Computed:    true,
						},
						"private_description": schema.StringAttribute{
							Description: "The private description of the service.",
							Computed:    true,
						},
						"parent_id": schema.StringAttribute{
							Description: "The service parent ID.",
							Computed:    true,
						},
						"current_incident_type": schema.StringAttribute{
							MarkdownDescription: "Enum: `\"major\"` `\"minor\"` `\"scheduled\"`\n  The service's current incident type.\n  The type of the (current) incident:\n" +
								"  - `minor` - A minor incident is currently taking place.\n" +
								"  - `major` - A major incident is currently taking place.\n" +
								"  - `scheduled` - A scheduled maintenance is currently taking place.",
							Computed: true,
						},
						"monitoring": schema.StringAttribute{
							MarkdownDescription: "Enum: `\"\"` `\"internal\"` `\"3rd_party\"` `\"webhook\"`\n  Monitoring types:\n" +
								"  - `\"\"` - No monitoring.\n" +
								"  - `internal` - StatusPal monitoring.\n" +
								"  - `3rd_party` - 3rd Party monitoring.\n" +
								"  - `webhook` - Incoming webhook monitoring.",
							Computed: true,
						},
						"webhook_monitoring_service": schema.StringAttribute{
							MarkdownDescription: "Enum: `\"status-cake\"` `\"uptime-robot\"` `\"custom-jsonpath\"`\n" +
								"  **Configure this field only if the `monitoring` is set to `webhook`.**\n" +
								"  Webhook Monitoring types:\n" +
								"  - `status-cake` - StatusCake monitoring service.\n" +
								"  - `internal` - UptimeRobot monitoring service.\n" +
								"  - `3rd_party` - Custom JSONPath.",
							Computed: true,
						},
						"webhook_custom_jsonpath_settings": schema.SingleNestedAttribute{
							MarkdownDescription: "The webhook monitoring service custom JSONPath settings.\n" +
								"  **Configure this field only if the `webhook_monitoring_service` is set to `custom-jsonpath`.**\n→ ",
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"jsonpath": schema.StringAttribute{
									MarkdownDescription: "The path in the JSON, e.g. `$.status`",
									Computed:            true,
								},
								"expected_result": schema.StringAttribute{
									MarkdownDescription: "The expected result in the JSON, e.g. `\"up\"`",
									Computed:            true,
								},
							},
						},
						"inbound_email_address": schema.StringAttribute{
							MarkdownDescription: "This is field is populated from `inbound_email_id`, if the `monitoring` is set to `3rd_party`.",
							Computed:            true,
						},
						"incoming_webhook_url": schema.StringAttribute{
							MarkdownDescription: "This is field is populated from `inbound_email_id`, if the `monitoring` is set to `webhook` and the `webhook_monitoring_service` is set.",
							Computed:            true,
						},
						"ping_url": schema.StringAttribute{
							Description: "We will send HTTP requests to this URL for monitoring every minute.",
							Computed:    true,
						},
						"incident_type": schema.StringAttribute{
							MarkdownDescription: "Enum: `\"major\"` `\"minor\"`\n  Sets the incident type to this value when an incident is created via monitoring.\n  The type of the (current) incident:\n" +
								"  - `minor` - A minor incident is currently taking place.\n" +
								"  - `major` - A major incident is currently taking place.",
							Computed: true,
						},
						"parent_incident_type": schema.StringAttribute{
							MarkdownDescription: "Enum: `\"major\"` `\"minor\"`\n  Sets the parent's service incident type to this value when an incident is created via monitoring.\n  The type of the (current) incident:\n" +
								"  - `minor` - A minor incident is currently taking place.\n" +
								"  - `major` - A major incident is currently taking place.",
							Computed: true,
						},
						"is_up": schema.BoolAttribute{
							Description: "Is the monitored service up?",
							Computed:    true,
						},
						"pause_monitoring_during_maintenances": schema.BoolAttribute{
							Description: "Pause the the service monitoring during maintenances?",
							Computed:    true,
						},
						"inbound_email_id": schema.StringAttribute{
							Description: "The inbound email ID.",
							Computed:    true,
						},
						"auto_incident": schema.BoolAttribute{
							Description: "Create an incident automatically when this service is down and close it if/when it comes back up.",
							Computed:    true,
						},
						"auto_notify": schema.BoolAttribute{
							Description: "Automatically notify all your subscribers about automatically created and closed incidents.",
							Computed:    true,
						},
						"children_ids": schema.ListAttribute{
							Description: "IDs of the service's children.",
							Computed:    true,
							ElementType: types.Int64Type,
						},
						"translations": schema.MapNestedAttribute{
							MarkdownDescription: "A translations object. For example:\n  ```terraform" + `
	{
		en = {
			name = "Your service"
			description = "This is your service's description..."
		}
		fr = {
			name = "Votre service"
			description = "Voici la description de votre service..."
		}
	}
` + "  ```\n→ ",
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "The name of the service.",
										Computed:    true,
									},
									"description": schema.StringAttribute{
										Description: "The description of the service.",
										Computed:    true,
									},
								},
							},
						},
						"private": schema.BoolAttribute{
							Description: "Private service?",
							Computed:    true,
						},
						"display_uptime_graph": schema.BoolAttribute{
							Description: "Display uptime graph?",
							Computed:    true,
						},
						"display_response_time_chart": schema.BoolAttribute{
							Description: "Display response time chart?",
							Computed:    true,
						},
						"order": schema.Int64Attribute{
							Description: "Service's position in the service list.",
							Computed:    true,
						},
						"inserted_at": schema.StringAttribute{
							Description: "Datetime at which the service was inserted.",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "Datetime at which the service was last updated.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *servicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Retrieve values from config
	var state servicesDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	statusPageSubdomain := state.StatusPageSubdomain.ValueString()
	services, err := d.client.GetServices(&statusPageSubdomain)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read StatusPal Services",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, service := range *services {
		// Create the translationData object dynamically
		translationData := make(servicesTranslationsModel)
		for lang, data := range service.Translations {
			translationData[lang] = servicesTranslationModel{
				Name:        types.StringValue(data.Name),
				Description: types.StringValue(data.Description),
			}
		}

		// Create the childrenIDsData list from service.ChildrenIDs
		childrenIDsData, diags := types.ListValueFrom(ctx, types.Int64Type, service.ChildrenIDs)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		serviceState := servicesModel{
			ID:                       types.StringValue(strconv.FormatInt(service.ID, 10)),
			Name:                     types.StringValue(service.Name),
			Description:              types.StringValue(service.Description),
			PrivateDescription:       types.StringValue(service.PrivateDescription),
			ParentID:                 types.StringValue(strconv.FormatInt(service.ParentID, 10)),
			CurrentIncidentType:      types.StringValue(service.CurrentIncidentType),
			Monitoring:               types.StringValue(service.Monitoring),
			WebhookMonitoringService: types.StringValue(service.WebhookMonitoringService),
			WebhookCustomJsonpathSettings: servicesWebhookCustomJsonpathSettingsModel{
				Jsonpath:       types.StringValue(service.WebhookCustomJsonpathSettings.Jsonpath),
				ExpectedResult: types.StringValue(service.WebhookCustomJsonpathSettings.ExpectedResult),
			},
			InboundEmailAddress:               types.StringValue(service.InboundEmailAddress),
			IncomingWebhookUrl:                types.StringValue(service.IncomingWebhookUrl),
			PingUrl:                           types.StringValue(service.PingUrl),
			IncidentType:                      types.StringValue(service.IncidentType),
			ParentIncidentType:                types.StringValue(service.ParentIncidentType),
			IsUp:                              types.BoolValue(service.IsUp),
			PauseMonitoringDuringMaintenances: types.BoolValue(service.PauseMonitoringDuringMaintenances),
			InboundEmailID:                    types.StringValue(service.InboundEmailID),
			AutoIncident:                      types.BoolValue(service.AutoIncident),
			AutoNotify:                        types.BoolValue(service.AutoNotify),
			ChildrenIDs:                       childrenIDsData,
			Translations:                      translationData,
			Private:                           types.BoolValue(service.Private),
			DisplayUptimeGraph:                types.BoolValue(service.DisplayUptimeGraph),
			DisplayResponseTimeChart:          types.BoolValue(service.DisplayResponseTimeChart),
			Order:                             types.Int64Value(service.Order),
			InsertedAt:                        types.StringValue(service.InsertedAt),
			UpdatedAt:                         types.StringValue(service.UpdatedAt),
		}

		state.Services = append(state.Services, serviceState)
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
func (d *servicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}
