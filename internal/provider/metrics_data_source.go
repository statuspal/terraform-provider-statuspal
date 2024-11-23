package provider

import (
	"context"
	"fmt"

	statuspal "terraform-provider-statuspal/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ datasource.DataSource              = &MetricsDataSource{}
	_ datasource.DataSourceWithConfigure = &MetricsDataSource{}
)

type MetricsDataSource struct {
	client *statuspal.Client
}

type queryMetrics struct {
	Before *string `tfsdk:"before"`
	After  *string `tfsdk:"after"`
	Limit  *int64  `tfsdk:"limit"`
}

type MetricsDataSourceModel struct {
	ID                  types.String  `tfsdk:"id"` // only for test case
	StatusPageSubdomain types.String  `tfsdk:"status_page_subdomain"`
	Query               types.Object  `tfsdk:"query"`
	Metrics             []metricModel `tfsdk:"metrics"`
}

func NewMetricsDataSource() datasource.DataSource {
	return &MetricsDataSource{}
}

func (d *MetricsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"query": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"before": schema.StringAttribute{
						Description: "Used as a cursor for pagination",
						Optional:    true,
					},
					"after": schema.StringAttribute{
						Description: "Used as a cursor for pagination",
						Optional:    true,
					},
					"limit": schema.Int64Attribute{
						Description: "Set the number of metrics to return in the response. This defaults to 20 items",
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.AtLeast(1),
							int64validator.AtMost(100),
						},
					},
				},
			},
		},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier attribute. Ignore it, only used in testing.",
				Computed:    true,
			},
			"status_page_subdomain": schema.StringAttribute{
				Description: "The status page subdomain of the services.",
				Required:    true,
			},
			"metrics": schema.ListNestedAttribute{
				Description: "The metrics",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "The unique identifier for the metric.",
							Required:    true,
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
							Computed:    true,
						},
						"unit": schema.StringAttribute{
							Description: "The unit of measurement for the metric.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "The type of the metric.",
							Computed:    true,
						},
						"enabled": schema.BoolAttribute{
							Description: "A flag indicating if the metric is enabled.",
							Computed:    true,
						},
						"visible": schema.BoolAttribute{
							Description: "A flag indicating if the metric is visible.",
							Computed:    true,
						},
						"remote_id": schema.StringAttribute{
							Description: "The remote ID for the metric.",
							Computed:    true,
						},
						"remote_name": schema.StringAttribute{
							Description: "The remote name for the metric.",
							Computed:    true,
						},
						"threshold": schema.Int64Attribute{
							Description: "The threshold value for the metric.",
							Computed:    true,
						},
						"featured_number": schema.StringAttribute{
							Description: "A featured number for the metric.",
							Computed:    true,
						},
						"integration_id": schema.StringAttribute{
							Description: "The integration ID related to the metric.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *MetricsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data MetricsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var query statuspal.MetricsQuery
	if !data.Query.IsNull() {
		var q queryMetrics
		resp.Diagnostics.Append(data.Query.As(ctx, &query, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		query.After = *q.After
		query.Before = *q.Before
		query.Limit = *q.Limit
	}

	metric, err := d.client.GetMetrics(data.StatusPageSubdomain.ValueString(), query)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get the metric, got error: %s", err))

		return
	}

	mapMetricsToDataSourceModel(metric, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *MetricsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MetricsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metrics"
}

func mapMetricsToDataSourceModel(metric *[]statuspal.Metric, data *MetricsDataSourceModel) {
	var metrics []metricModel

	for _, m := range *metric {
		metrics = append(metrics, metricModel{
			ID:              types.Int64Value(m.ID),
			Status:          types.StringValue(m.Status),
			LatestEntryTime: types.Int64Value(m.LatestEntryTime),
			Order:           types.Int64Value(m.Order),
			Title:           types.StringValue(m.Title),
			Unit:            types.StringValue(m.Unit),
			Type:            types.StringValue(string(m.Type)),
			Enabled:         types.BoolValue(m.Enabled),
			Visible:         types.BoolValue(m.Visible),
			RemoteID:        types.StringValue(m.RemoteID),
			RemoteName:      types.StringValue(m.RemoteName),
			Threshold:       types.Int64Value(m.Threshold),
			FeaturedNumber:  types.StringValue(string(m.FeaturedNumber)),
			IntegrationID:   types.StringValue(m.IntegrationID),
		})
	}

	data.Metrics = metrics
}
