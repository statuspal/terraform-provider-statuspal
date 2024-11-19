package provider

import (
	"context"
	"fmt"

	statuspal "terraform-provider-statuspal/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &MetricResource{}
	_ resource.ResourceWithImportState = &MetricResource{}
)

func NewMetricResource() resource.Resource {
	return &MetricResource{}
}

// MetricResource defines the resource implementation.
type MetricResource struct {
	client *statuspal.Client
}

// MetricType represents the type of metric (either "up" or "rt").
type MetricType string

const (
	UptimeMetric       MetricType = "up"
	ResponseTimeMetric MetricType = "rt"
)

// FeaturedNumber represents the type of number to display for the metric.
type FeaturedNumber string

const (
	AvgFeatured  FeaturedNumber = "avg"
	MaxFeatured  FeaturedNumber = "max"
	LastFeatured FeaturedNumber = "last"
)

type metricModel struct {
	ID              types.Int64  `tfsdk:"id"`
	Title           types.String `tfsdk:"title"`
	Unit            types.String `tfsdk:"unit"`
	Type            types.String `tfsdk:"type"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	Visible         types.Bool   `tfsdk:"visible"`
	RemoteID        types.String `tfsdk:"remote_id"`
	RemoteName      types.String `tfsdk:"remote_name"`
	Status          types.String `tfsdk:"status"`
	LatestEntryTime types.Int64  `tfsdk:"latest_entry_time"`
	Threshold       types.Int64  `tfsdk:"threshold"`
	FeaturedNumber  types.String `tfsdk:"featured_number"`
	Order           types.Int64  `tfsdk:"order"`
	IntegrationID   types.String `tfsdk:"integration_id"`
}

// MetricResourceModel represents the data model for a metric resource.
//
// https://www.statuspal.io/api-docs#tag/Metrics/operation/addMetric
type MetricResourceModel struct {
	ID                  types.String `tfsdk:"id"` // only for test case
	StatusPageSubdomain types.String `tfsdk:"status_page_subdomain"`
	Metric              metricModel  `tfsdk:"metric"`
}

func (r *MetricResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metric"
}

func (r *MetricResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier attribute. Ignore it, only used in testing.",
				Computed:    true,
			},
			"status_page_subdomain": schema.StringAttribute{
				Description: "The status page subdomain of the services.",
				Required:    true,
			},
			"metric": schema.SingleNestedAttribute{
				Description: "The metric.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description: "The unique identifier for the metric.",
						Computed:    true,
					},
					"status": schema.StringAttribute{
						Description: "The status of the metric.",
						Computed:    true,
					},
					"latest_entry_time": schema.Int64Attribute{
						Description: "The timestamp for the latest entry of the metric.",
						Computed:    true,
					},
					"order": schema.Int64Attribute{
						Description: "The order of the metric in the system.",
						Computed:    true,
					},
					"title": schema.StringAttribute{
						Description: "The title of the metric.",
						Required:    true,
					},
					"unit": schema.StringAttribute{
						Description: "The unit of measurement for the metric.",
						Required:    true,
					},
					"type": schema.StringAttribute{
						Description: "The type of the metric.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf(string(UptimeMetric), string(ResponseTimeMetric)),
						},
					},
					"enabled": schema.BoolAttribute{
						Description: "A flag indicating if the metric is enabled.",
						Optional:    true,
						Computed:    true,
					},
					"visible": schema.BoolAttribute{
						Description: "A flag indicating if the metric is visible.",
						Optional:    true,
						Computed:    true,
					},
					"remote_id": schema.StringAttribute{
						Description: "The remote ID for the metric.",
						Optional:    true,
						Computed:    true,
					},
					"remote_name": schema.StringAttribute{
						Description: "The remote name for the metric.",
						Optional:    true,
						Computed:    true,
					},
					"threshold": schema.Int64Attribute{
						Description: "The threshold value for the metric.",
						Optional:    true,
						Computed:    true,
					},
					"featured_number": schema.StringAttribute{
						Description: "A featured number for the metric.",
						Optional:    true,
						Computed:    true,
						Validators: []validator.String{
							stringvalidator.OneOf(string(AvgFeatured), string(LastFeatured), string(MaxFeatured)),
						},
					},
					"integration_id": schema.StringAttribute{
						Description: "The integration ID related to the metric.",
						Optional:    true,
						Computed:    true,
					},
				},
			},
		},
	}
}

func (r *MetricResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// https://www.statuspal.io/api-docs#tag/Metrics/operation/addMetric
func (r *MetricResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MetricResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var model statuspal.Metric
	mapResourceModelToMetric(&model, &data)

	metric, err := r.client.CreateMetric(data.StatusPageSubdomain.ValueString(), &model)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create the metric, got error: %s", err))

		return
	}

	data.ID = types.StringValue("placeholder") // only for test case

	mapMetricToResourceModel(metric, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// https://www.statuspal.io/api-docs#tag/Metrics/operation/getMetric
func (r *MetricResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MetricResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	metric, err := r.client.GetMetric(data.Metric.ID.ValueInt64(), data.StatusPageSubdomain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get the metric, got error: %s", err))

		return
	}

	mapMetricToResourceModel(metric, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// https://www.statuspal.io/api-docs#tag/Metrics/operation/updateMetric
func (r *MetricResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data MetricResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Metric.ID.ValueInt64()
	subdomain := data.StatusPageSubdomain.ValueString()

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var model statuspal.Metric
	mapResourceModelToMetric(&model, &data)

	metric, err := r.client.UpdateMetric(id, subdomain, &model)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update the metric, got error: %s", err))

		return
	}

	mapMetricToResourceModel(metric, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// https://www.statuspal.io/api-docs#tag/Metrics/operation/deleteMetric
func (r *MetricResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MetricResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteMetric(data.Metric.ID.ValueInt64(), data.StatusPageSubdomain.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete the metric, got error: %s", err))

		return
	}
}

func (r *MetricResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapMetricToResourceModel(metric *statuspal.Metric, data *MetricResourceModel) {
	data.Metric.ID = types.Int64Value(metric.ID)
	data.Metric.Title = types.StringValue(metric.Title)
	data.Metric.Unit = types.StringValue(metric.Unit)
	data.Metric.Type = types.StringValue(string(metric.Type))
	data.Metric.Enabled = types.BoolValue(metric.Enabled)
	data.Metric.Visible = types.BoolValue(metric.Visible)
	data.Metric.RemoteID = types.StringValue(metric.RemoteID)
	data.Metric.RemoteName = types.StringValue(metric.RemoteName)
	data.Metric.Status = types.StringValue(metric.Status)
	data.Metric.LatestEntryTime = types.Int64Value(metric.LatestEntryTime)
	data.Metric.Threshold = types.Int64Value(metric.Threshold)
	data.Metric.FeaturedNumber = types.StringValue(string(metric.FeaturedNumber))
	data.Metric.Order = types.Int64Value(metric.Order)
	data.Metric.IntegrationID = types.StringValue(metric.IntegrationID)
}

func mapResourceModelToMetric(metric *statuspal.Metric, data *MetricResourceModel) {
	metric.ID = data.Metric.ID.ValueInt64()
	metric.Title = data.Metric.Title.ValueString()
	metric.Unit = data.Metric.Unit.ValueString()
	metric.Type = statuspal.MetricType(data.Metric.Type.ValueString())
	metric.Enabled = data.Metric.Enabled.ValueBool()
	metric.Visible = data.Metric.Visible.ValueBool()
	metric.RemoteID = data.Metric.RemoteID.ValueString()
	metric.RemoteName = data.Metric.RemoteName.ValueString()
	metric.Status = data.Metric.Status.ValueString()
	metric.LatestEntryTime = data.Metric.LatestEntryTime.ValueInt64()
	metric.Threshold = data.Metric.Threshold.ValueInt64()
	metric.FeaturedNumber = statuspal.FeaturedNumber(data.Metric.FeaturedNumber.ValueString())
	metric.Order = data.Metric.Order.ValueInt64()
	metric.IntegrationID = data.Metric.IntegrationID.ValueString()
}