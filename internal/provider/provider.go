package provider

import (
	"context"
	"os"
	"strings"

	statuspal "terraform-provider-statuspal/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	ApiKey types.String `tfsdk:"api_key"`
	Region types.String `tfsdk:"region"`
}

type statuspalProviderDevModel struct {
	ApiKey types.String `tfsdk:"api_key"`
}

type statuspalProviderTestModel struct {
	TestUrl types.String `tfsdk:"test_url"`
}

// Metadata returns the provider type name.
func (p *statuspalProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "statuspal"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *statuspalProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	env := os.Getenv("TF_ENV")

	attributes := map[string]schema.Attribute{}

	if env == "DEV" || env != "TEST" {
		attributes["api_key"] = schema.StringAttribute{
			MarkdownDescription: "Your StatusPal User or Organization API Key. May also be provided via `STATUSPAL_API_KEY` environment variable.",
			Optional:            true,
			Sensitive:           true,
		}

		if env != "DEV" {
			attributes["region"] = schema.StringAttribute{
				MarkdownDescription: "StatusPal API Region, it can be \"US\" and \"EU\". May also be provided via `STATUSPAL_REGION` environment variable.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("EU", "US"),
				},
			}
		}
	} else {
		attributes["test_url"] = schema.StringAttribute{
			Required: true,
		}
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Interact with [StatusPal](https://www.statuspal.io).",
		Attributes:          attributes,
	}
}

// Configure prepares a StatusPal API client for data sources and resources.
func (p *statuspalProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring StatusPal client")

	env := os.Getenv("TF_ENV")

	var api_key string
	var region string
	var test_url string

	if env == "DEV" {
		// Retrieve provider data from configuration
		var config statuspalProviderDevModel
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
			return
		}

		// Default values to environment variables, but override
		// with Terraform configuration value if set.

		api_key = os.Getenv("STATUSPAL_API_KEY")

		if !config.ApiKey.IsNull() {
			api_key = config.ApiKey.ValueString()
		}

		// If any of the expected configurations are missing, return
		// errors with provider-specific guidance.

		if api_key == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("api_key"),
				"Missing StatusPal API Key",
				"The provider cannot create the StatusPal API client as there is a missing or empty value for the StatusPal API key. "+
					"Set the api key value in the configuration or use the STATUSPAL_API_KEY environment variable. "+
					"If either is already set, ensure the value is not empty.",
			)
			return
		}

		ctx = tflog.SetField(ctx, "api_key", api_key)
	} else if env == "TEST" {
		// Retrieve provider data from configuration
		var config statuspalProviderTestModel
		diags := req.Config.Get(ctx, &config)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		test_url = config.TestUrl.ValueString()
	} else {
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
			return
		}

		if config.Region.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("region"),
				"Unknown StatusPal API Region",
				"The provider cannot create the StatusPal API client as there is an unknown configuration value for the StatusPal API region. "+
					"Either target apply the source of the value first, set the value statically in the configuration, or use the STATUSPAL_REGION environment variable.",
			)
			return
		}

		// Default values to environment variables, but override
		// with Terraform configuration value if set.

		api_key = os.Getenv("STATUSPAL_API_KEY")
		region = os.Getenv("STATUSPAL_REGION")

		if !config.ApiKey.IsNull() {
			api_key = config.ApiKey.ValueString()
		}

		if !config.Region.IsNull() {
			region = config.Region.ValueString()
		}

		region = strings.ToUpper(region)

		// If any of the expected configurations are missing, return
		// errors with provider-specific guidance.

		if api_key == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("api_key"),
				"Missing StatusPal API Key",
				"The provider cannot create the StatusPal API client as there is a missing or empty value for the StatusPal API key. "+
					"Set the api key value in the configuration or use the STATUSPAL_API_KEY environment variable. "+
					"If either is already set, ensure the value is not empty.",
			)
			return
		}

		if region != "EU" && region != "US" {
			resp.Diagnostics.AddAttributeError(
				path.Root("region"),
				"Missing or Invalid StatusPal API Region",
				"The provider cannot create the StatusPal API client as there is a missing, empty or invalid value for the StatusPal API region. "+
					"Set the region value in the configuration or use the STATUSPAL_REGION environment variable. "+
					`If either is already set, ensure the value is not empty and it can be only "EU" or "US".`,
			)
			return
		}

		ctx = tflog.SetField(ctx, "api_key", api_key)
		ctx = tflog.SetField(ctx, "region", region)
	}

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
