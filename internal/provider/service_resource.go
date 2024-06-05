package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	statuspal "terraform-provider-statuspal/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &serviceResource{}
	_ resource.ResourceWithConfigure   = &serviceResource{}
	_ resource.ResourceWithImportState = &serviceResource{}
)

// NewServiceResource is a helper function to simplify the provider implementation.
func NewServiceResource() resource.Resource {
	return &serviceResource{}
}

// serviceResource is the resource implementation.
type serviceResource struct {
	client *statuspal.Client
}

// serviceResourceModel maps the resource schema data.
type serviceResourceModel struct {
	ID                  types.String `tfsdk:"id"` // only for test case
	StatusPageSubdomain types.String `tfsdk:"status_page_subdomain"`
	Service             serviceModel `tfsdk:"service"`
}

// serviceModel maps service schema data.
type serviceModel struct {
	ID                                types.String `tfsdk:"id"`
	Name                              types.String `tfsdk:"name"`
	Description                       types.String `tfsdk:"description"`
	PrivateDescription                types.String `tfsdk:"private_description"`
	CurrentIncidentType               types.String `tfsdk:"current_incident_type"`
	Monitoring                        types.String `tfsdk:"monitoring"`
	PingUrl                           types.String `tfsdk:"ping_url"`
	IncidentType                      types.String `tfsdk:"incident_type"`
	ParentIncidentType                types.String `tfsdk:"parent_incident_type"`
	IsUp                              types.Bool   `tfsdk:"is_up"`
	PauseMonitoringDuringMaintenances types.Bool   `tfsdk:"pause_monitoring_during_maintenances"`
	InboundEmailID                    types.String `tfsdk:"inbound_email_id"`
	AutoIncident                      types.Bool   `tfsdk:"auto_incident"`
	AutoNotify                        types.Bool   `tfsdk:"auto_notify"`
	ChildrenIDs                       types.List   `tfsdk:"children_ids"`
	Translations                      types.Map    `tfsdk:"translations"`
	Private                           types.Bool   `tfsdk:"private"`
	DisplayUptimeGraph                types.Bool   `tfsdk:"display_uptime_graph"`
	DisplayResponseTimeChart          types.Bool   `tfsdk:"display_response_time_chart"`
	Order                             types.Int64  `tfsdk:"order"`
	InsertedAt                        types.String `tfsdk:"inserted_at"`
	UpdatedAt                         types.String `tfsdk:"updated_at"`
}

type serviceTranslationsModel map[string]serviceTranslationModel

type serviceTranslationModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

// Metadata returns the resource type name.
func (r *serviceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

// Schema defines the schema for the resource.
func (r *serviceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a service of the status page.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier attribute. Ignore it, only used in testing.",
				Computed:    true,
			},
			"status_page_subdomain": schema.StringAttribute{
				Description: "The status page's subdomain where the service belong.",
				Required:    true,
			},
			"service": schema.SingleNestedAttribute{
				Description: "The service.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The ID of the service.",
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						Description: "The name of the service.",
						Required:    true,
					},
					"description": schema.StringAttribute{
						Description: "The description of the service.",
						Optional:    true,
						Computed:    true,
					},
					"private_description": schema.StringAttribute{
						Description: "The private description of the service.",
						Optional:    true,
						Computed:    true,
					},
					"current_incident_type": schema.StringAttribute{
						MarkdownDescription: "Enum: `\"major\"` `\"minor\"` `\"scheduled\"`\n  The type of the (current) incident:\n" +
							"  - `major` - A minor incident is currently taking place.\n" +
							"  - `minor` - A major incident is currently taking place.\n" +
							"  - `scheduled` - A scheduled maintenance is currently taking place.",
						Optional: true,
						Computed: true,
					},
					"monitoring": schema.StringAttribute{
						MarkdownDescription: "Enum: `null` `\"internal\"` `\"3rd_party\"`\n  Monitoring types:\n" +
							"  - `major` - No monitoring.\n" +
							"  - `internal` - StatusPal monitoring.\n" +
							"  - `3rd_party` - 3rd Party monitoring.",
						Optional: true,
						Computed: true,
					},
					"ping_url": schema.StringAttribute{
						Description: "We will send HTTP requests to this URL for monitoring every minute.",
						Optional:    true,
						Computed:    true,
					},
					"incident_type": schema.StringAttribute{
						MarkdownDescription: "Enum: `\"major\"` `\"minor\"`\n  Sets the incident type to this value when an incident is created via monitoring.\n  The type of the (current) incident:\n" +
							"  - `major` - A minor incident is currently taking place.\n" +
							"  - `minor` - A major incident is currently taking place.",
						Optional: true,
						Computed: true,
					},
					"parent_incident_type": schema.StringAttribute{
						MarkdownDescription: "Enum: `\"major\"` `\"minor\"`\n  Sets the parent's service incident type to this value when an incident is created via monitoring.\n  The type of the (current) incident:\n" +
							"  - `major` - A minor incident is currently taking place.\n" +
							"  - `minor` - A major incident is currently taking place.",
						Optional: true,
						Computed: true,
					},
					"is_up": schema.BoolAttribute{
						Description: "Is the monitored service up?",
						Optional:    true,
						Computed:    true,
					},
					"pause_monitoring_during_maintenances": schema.BoolAttribute{
						Description: "Pause the the service monitoring during maintenances?",
						Optional:    true,
						Computed:    true,
					},
					"inbound_email_id": schema.StringAttribute{
						Description: "The inbound email ID.",
						Computed:    true,
					},
					"auto_incident": schema.BoolAttribute{
						Description: "Create an incident automatically when this service is down and close it if/when it comes back up.",
						Optional:    true,
						Computed:    true,
					},
					"auto_notify": schema.BoolAttribute{
						Description: "Automatically notify all your subscribers about automatically created and closed incidents.",
						Optional:    true,
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
						Optional: true,
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of the service.",
									Required:    true,
								},
								"description": schema.StringAttribute{
									Description: "The description of the service.",
									Required:    true,
								},
							},
						},
					},
					"private": schema.BoolAttribute{
						Description: "Private service?",
						Optional:    true,
						Computed:    true,
					},
					"display_uptime_graph": schema.BoolAttribute{
						Description: "Display uptime graph?",
						Optional:    true,
						Computed:    true,
					},
					"display_response_time_chart": schema.BoolAttribute{
						Description: "Display response time chart?",
						Optional:    true,
						Computed:    true,
					},
					"order": schema.Int64Attribute{
						Description: "Service's position in the service list.",
						Optional:    true,
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
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *serviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan serviceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	service := mapServiceModelToRequestBody(&ctx, &plan.Service, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new service
	statusPageSubdomain := plan.StatusPageSubdomain.ValueString()
	newService, err := r.client.CreateService(service, &statusPageSubdomain)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating StatusPal Service",
			"Could not create service, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	newServiceModel := mapResponseToServiceModel(&ctx, newService, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Service = *newServiceModel
	plan.ID = types.StringValue("placeholder") // only for test case

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *serviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state serviceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed service value from StatusPal
	statusPageSubdomain := state.StatusPageSubdomain.ValueString()
	serviceID := state.Service.ID.ValueString()
	service, err := r.client.GetService(&statusPageSubdomain, &serviceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading StatusPal Service",
			"Could not read service ID "+serviceID+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	serviceModel := mapResponseToServiceModel(&ctx, service, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Service = *serviceModel
	state.ID = types.StringValue("placeholder") // only for test case

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *serviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan serviceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	service := mapServiceModelToRequestBody(&ctx, &plan.Service, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing service
	statusPageSubdomain := plan.StatusPageSubdomain.ValueString()
	serviceID := plan.Service.ID.ValueString()
	updatedService, err := r.client.UpdateService(service, &statusPageSubdomain, &serviceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating StatusPal Service",
			"Could not Update service, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	updatedServiceModel := mapResponseToServiceModel(&ctx, updatedService, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Service = *updatedServiceModel
	plan.ID = types.StringValue("placeholder") // only for test case

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *serviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state serviceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	statusPageSubdomain := state.StatusPageSubdomain.ValueString()
	serviceID := state.Service.ID.ValueString()
	err := r.client.DeleteService(&statusPageSubdomain, &serviceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting StatusPal Service",
			"Could not delete service, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *serviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Split the ID based on the delimiter used during import
	parts := strings.Split(req.ID, " ")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected StatusPal Service Import Identifier",
			`Expected StatusPal service import identifier with format: "<status_page_subdomain> <service_id>"`,
		)
		return
	}

	req.ID = parts[0]
	resource.ImportStatePassthroughID(ctx, path.Root("status_page_subdomain"), req, resp)
	req.ID = parts[1]
	resource.ImportStatePassthroughID(ctx, path.Root("service").AtName("id"), req, resp)
}

// Configure adds the provider configured client to the resource.
func (r *serviceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func mapServiceModelToRequestBody(ctx *context.Context, service *serviceModel, diagnostics *diag.Diagnostics) *statuspal.Service {
	// Create the translationData object dynamically
	translationData := make(statuspal.ServiceTranslations)
	if !service.Translations.IsNull() && !service.Translations.IsUnknown() {
		translations := make(serviceTranslationsModel, len(service.Translations.Elements()))
		diags := service.Translations.ElementsAs(*ctx, &translations, false)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return nil
		}

		for lang, data := range translations {
			translationData[lang] = statuspal.ServiceTranslation{
				Name:        data.Name.ValueString(),
				Description: data.Description.ValueString(),
			}
		}
	}

	return &statuspal.Service{
		Name:                              service.Name.ValueString(),
		Description:                       service.Description.ValueString(),
		PrivateDescription:                service.PrivateDescription.ValueString(),
		CurrentIncidentType:               service.CurrentIncidentType.ValueString(),
		Monitoring:                        service.Monitoring.ValueString(),
		PingUrl:                           service.PingUrl.ValueString(),
		IncidentType:                      service.IncidentType.ValueString(),
		ParentIncidentType:                service.ParentIncidentType.ValueString(),
		IsUp:                              service.IsUp.ValueBool(),
		PauseMonitoringDuringMaintenances: service.PauseMonitoringDuringMaintenances.ValueBool(),
		InboundEmailID:                    service.InboundEmailID.ValueString(),
		AutoIncident:                      service.AutoIncident.ValueBool(),
		AutoNotify:                        service.AutoNotify.ValueBool(),
		Translations:                      translationData,
		Private:                           service.Private.ValueBool(),
		DisplayUptimeGraph:                service.DisplayUptimeGraph.ValueBool(),
		DisplayResponseTimeChart:          service.DisplayResponseTimeChart.ValueBool(),
		Order:                             service.Order.ValueInt64(),
	}
}

func mapResponseToServiceModel(ctx *context.Context, service *statuspal.Service, diagnostics *diag.Diagnostics) *serviceModel {
	// Define the translation object schema
	translationSchema := map[string]attr.Type{
		"name":        types.StringType,
		"description": types.StringType,
	}
	// Create the translationData object dynamically
	translationData := make(map[string]attr.Value)
	for lang, data := range service.Translations {
		translationObject, diags := types.ObjectValue(
			translationSchema,
			map[string]attr.Value{
				"name":        types.StringValue(data.Name),
				"description": types.StringValue(data.Description),
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

	// Create the childrenIDsData list from service.ChildrenIDs
	childrenIDsData, diags := types.ListValueFrom(*ctx, types.Int64Type, service.ChildrenIDs)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return nil
	}

	return &serviceModel{
		ID:                                types.StringValue(strconv.FormatInt(service.ID, 10)),
		Name:                              types.StringValue(service.Name),
		Description:                       types.StringValue(service.Description),
		PrivateDescription:                types.StringValue(service.PrivateDescription),
		CurrentIncidentType:               types.StringValue(service.CurrentIncidentType),
		Monitoring:                        types.StringValue(service.Monitoring),
		PingUrl:                           types.StringValue(service.PingUrl),
		IncidentType:                      types.StringValue(service.IncidentType),
		ParentIncidentType:                types.StringValue(service.ParentIncidentType),
		IsUp:                              types.BoolValue(service.IsUp),
		PauseMonitoringDuringMaintenances: types.BoolValue(service.PauseMonitoringDuringMaintenances),
		InboundEmailID:                    types.StringValue(service.InboundEmailID),
		AutoIncident:                      types.BoolValue(service.AutoIncident),
		AutoNotify:                        types.BoolValue(service.AutoNotify),
		ChildrenIDs:                       childrenIDsData,
		Translations:                      translations,
		Private:                           types.BoolValue(service.Private),
		DisplayUptimeGraph:                types.BoolValue(service.DisplayUptimeGraph),
		DisplayResponseTimeChart:          types.BoolValue(service.DisplayResponseTimeChart),
		Order:                             types.Int64Value(service.Order),
		InsertedAt:                        types.StringValue(service.InsertedAt),
		UpdatedAt:                         types.StringValue(service.UpdatedAt),
	}
}
