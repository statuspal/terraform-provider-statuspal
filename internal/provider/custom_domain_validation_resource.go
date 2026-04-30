package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"

	statuspal "terraform-provider-statuspal/internal/client"
)

const (
	domainValidationPollInterval = 10 * time.Second
	domainValidationDefaultTimeout = 30 * time.Minute
)

var (
	_ resource.Resource              = &customDomainValidationResource{}
	_ resource.ResourceWithConfigure = &customDomainValidationResource{}
)

func NewCustomDomainValidationResource() resource.Resource {
	return &customDomainValidationResource{}
}

type customDomainValidationResource struct {
	client *statuspal.Client
}

type customDomainValidationResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	OrganizationID      types.String `tfsdk:"organization_id"`
	StatusPageSubdomain types.String `tfsdk:"status_page_subdomain"`
	TimeoutSeconds      types.Int64  `tfsdk:"timeout_seconds"`
}

func (r *customDomainValidationResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_custom_domain_validation"
}

func (r *customDomainValidationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Waiter resource that blocks until a status page's custom domain reaches the \"active\" state. " +
			"Use this after creating DNS records to ensure the domain is verified before proceeding. " +
			"Destroying this resource is a no-op — it only removes the waiter from Terraform state.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier attribute.",
				Computed:    true,
			},
			"organization_id": schema.StringAttribute{
				Description: "The organization ID that owns the status page.",
				Required:    true,
			},
			"status_page_subdomain": schema.StringAttribute{
				Description: "The subdomain of the status page whose custom domain should be validated.",
				Required:    true,
			},
			"timeout_seconds": schema.Int64Attribute{
				Description: "Maximum seconds to wait for the domain to become active. Defaults to 1800 (30 minutes).",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1800),
			},
		},
	}
}

func (r *customDomainValidationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan customDomainValidationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.pollUntilActive(ctx, plan.OrganizationID.ValueString(), plan.StatusPageSubdomain.ValueString(), plan.TimeoutSeconds.ValueInt64()); err != nil {
		resp.Diagnostics.AddError("Custom domain validation failed", err.Error())
		return
	}

	plan.ID = types.StringValue("placeholder")
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *customDomainValidationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state customDomainValidationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID := state.OrganizationID.ValueString()
	subdomain := state.StatusPageSubdomain.ValueString()
	statusPage, err := r.client.GetStatusPage(&orgID, &subdomain)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading status page for domain validation",
			"Could not read status page "+subdomain+": "+err.Error(),
		)
		return
	}

	// If the domain config is gone or disabled, remove this resource from state
	// so Terraform knows to re-create it on next apply.
	if statusPage.DomainConfig == nil || statusPage.DomainConfig.Status == nil ||
		*statusPage.DomainConfig.Status == "disabled" {
		resp.State.RemoveResource(ctx)
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update re-polls if the subdomain or timeout changes.
func (r *customDomainValidationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan customDomainValidationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.pollUntilActive(ctx, plan.OrganizationID.ValueString(), plan.StatusPageSubdomain.ValueString(), plan.TimeoutSeconds.ValueInt64()); err != nil {
		resp.Diagnostics.AddError("Custom domain validation failed", err.Error())
		return
	}

	plan.ID = types.StringValue("placeholder")
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete is a no-op: removing the waiter from config does not affect the domain.
func (r *customDomainValidationResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

func (r *customDomainValidationResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*statuspal.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *statuspal.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *customDomainValidationResource) pollUntilActive(ctx context.Context, orgID, subdomain string, timeoutSeconds int64) error {
	timeout := time.Duration(timeoutSeconds) * time.Second
	deadline := time.Now().Add(timeout)

	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("timed out after %ds waiting for custom domain on status page %q to become active", timeoutSeconds, subdomain)
		}

		// Respect context cancellation (e.g. user presses Ctrl+C).
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while waiting for custom domain on status page %q", subdomain)
		default:
		}

		statusPage, err := r.client.GetStatusPage(&orgID, &subdomain)
		if err != nil {
			return fmt.Errorf("error polling status page %q: %w", subdomain, err)
		}

		if statusPage.DomainConfig == nil || statusPage.DomainConfig.Status == nil {
			return fmt.Errorf("status page %q has no domain_config; ensure domain_config is set before using this resource", subdomain)
		}

		switch *statusPage.DomainConfig.Status {
		case "active":
			return nil
		case "failed_to_configure":
			errMsg := ""
			if statusPage.DomainConfig.Error != nil {
				errMsg = *statusPage.DomainConfig.Error
			}
			return fmt.Errorf("custom domain on status page %q failed to configure: %s", subdomain, errMsg)
		case "disabled":
			return fmt.Errorf("custom domain on status page %q is disabled; set domain_config before using this resource", subdomain)
		}

		// Status is "configuring" — wait and retry.
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while waiting for custom domain on status page %q", subdomain)
		case <-time.After(domainValidationPollInterval):
		}
	}
}
