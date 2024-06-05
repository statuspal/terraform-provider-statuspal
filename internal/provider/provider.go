package provider

import (
	"context"
	"os"
	"strings"

	statuspal "terraform-provider-statuspal/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider              = &statuspalProvider{}
	_ provider.ProviderWithFunctions = &statuspalProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &statuspalProvider{
			version: version,
		}
	}
}

// statuspalProvider is the provider implementation.
type statuspalProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// statuspalProviderModel maps provider schema data to a Go type.
type statuspalProviderModel struct {
	ApiKey  types.String `tfsdk:"api_key"`
	Region  types.String `tfsdk:"region"`
	TestUrl types.String `tfsdk:"test_url"`
}

// Metadata returns the provider type name.
func (p *statuspalProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "statuspal"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *statuspalProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with StatusPal.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: "Your StatusPal User or Organization API Key. May also be provided via STATUSPAL_API_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"region": schema.StringAttribute{
				Description: `StatusPal API Region, it can be "US" and "EU". May also be provided via STATUSPAL_REGION environment variable.`,
				Optional:    true,
			},
			"test_url": schema.StringAttribute{
				Description: "Ignore this attribute, it's only used in testing.",
				Optional:    true,
			},
		},
	}
}

// Configure prepares a StatusPal API client for data sources and resources.
func (p *statuspalProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring StatusPal client")

	// Retrieve provider data from configuration
	var config statuspalProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown StatusPal API Key",
			"The provider cannot create the StatusPal API client as there is an unknown configuration value for the StatusPal API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the STATUSPAL_API_KEY environment variable.",
		)
	}

	if config.Region.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("region"),
			"Unknown StatusPal API Region",
			"The provider cannot create the StatusPal API client as there is an unknown configuration value for the StatusPal API region. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the STATUSPAL_REGION environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	api_key := os.Getenv("STATUSPAL_API_KEY")
	region := os.Getenv("STATUSPAL_REGION")

	if !config.ApiKey.IsNull() {
		api_key = config.ApiKey.ValueString()
	}

	if !config.Region.IsNull() {
		region = config.Region.ValueString()
	}

	region = strings.ToLower(region)

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if api_key == "" && region != "test" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing StatusPal API Key",
			"The provider cannot create the StatusPal API client as there is a missing or empty value for the StatusPal API key. "+
				"Set the api key value in the configuration or use the STATUSPAL_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if region != "eu" && region != "us" && region != "dev" && region != "test" {
		resp.Diagnostics.AddAttributeError(
			path.Root("region"),
			"Missing or Invalid StatusPal API Region",
			"The provider cannot create the StatusPal API client as there is a missing or empty or invalid value for the StatusPal API region. "+
				"Set the region value in the configuration or use the STATUSPAL_REGION environment variable. "+
				`If either is already set, ensure the value is not empty and it can be only "EU" or "US".`,
		)
	}

	if region == "test" && (config.TestUrl.IsUnknown() || config.TestUrl.IsNull() || config.TestUrl.ValueString() == "") {
		resp.Diagnostics.AddAttributeError(
			path.Root("test_url"),
			"Missing or Invalid TestUrl",
			"When you are testing always has to provider test_url in the provider config. You can create a mockServer for testing and use that URL.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	test_url := config.TestUrl.ValueString()

	ctx = tflog.SetField(ctx, "api_key", api_key)
	ctx = tflog.SetField(ctx, "region", region)
	ctx = tflog.SetField(ctx, "test_url", test_url)

	tflog.Debug(ctx, "Creating StatusPal client")

	// Create a new StatusPal client using the configuration values
	client, err := statuspal.NewClient(&api_key, &region, &test_url)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create StatusPal API Client",
			"An unexpected error occurred when creating the StatusPal API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"StatusPal Client Error: "+err.Error(),
		)
		return
	}

	// Make the StatusPal client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured StatusPal client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *statuspalProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewStatusPagesDataSource,
		NewServicesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *statuspalProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewStatusPageResource,
		NewServiceResource,
	}
}

func (p *statuspalProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		NewDemoSumFunction,
	}
}
