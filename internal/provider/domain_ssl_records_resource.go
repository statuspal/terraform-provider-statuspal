package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	statuspal "terraform-provider-statuspal/internal/client"
)

const (
	sslRecordsPollInterval   = 10 * time.Second
	sslRecordsDefaultTimeout = 5 * time.Minute
)

var (
	_ resource.Resource              = &domainSslRecordsResource{}
	_ resource.ResourceWithConfigure = &domainSslRecordsResource{}
)

func NewDomainSslRecordsResource() resource.Resource {
	return &domainSslRecordsResource{}
}

type domainSslRecordsResource struct {
	client *statuspal.Client
}

type domainSslRecordsResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	OrganizationID      types.String `tfsdk:"organization_id"`
	StatusPageSubdomain types.String `tfsdk:"status_page_subdomain"`
	TimeoutSeconds      types.Int64  `tfsdk:"timeout_seconds"`
	CertificateTxtName  types.String `tfsdk:"certificate_txt_name"`
	CertificateTxtValue types.String `tfsdk:"certificate_txt_value"`
}

func (r *domainSslRecordsResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_domain_ssl_records"
}

func (r *domainSslRecordsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Waiter resource that polls a status page's domain_config until the SSL certificate " +
			"challenge DNS records become available, then exposes them as computed attributes. " +
			"Use this between creating the CNAME routing record and the TXT certificate record " +
			"to enable a single-apply custom domain flow. Destroying this resource is a no-op.",
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
				Description: "The subdomain of the status page.",
				Required:    true,
			},
			"timeout_seconds": schema.Int64Attribute{
				Description: "Maximum seconds to wait for the SSL certificate records to appear. Defaults to 300 (5 minutes).",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(300),
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"certificate_txt_name": schema.StringAttribute{
				Description: "The DNS name for the TXT record required to issue the SSL certificate.",
				Computed:    true,
			},
			"certificate_txt_value": schema.StringAttribute{
				Description: "The DNS value for the TXT record required to issue the SSL certificate.",
				Computed:    true,
			},
		},
	}
}

func (r *domainSslRecordsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan domainSslRecordsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name, value, err := r.pollUntilCertRecordsReady(
		ctx,
		plan.OrganizationID.ValueString(),
		plan.StatusPageSubdomain.ValueString(),
		plan.TimeoutSeconds.ValueInt64(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Failed waiting for SSL certificate records", err.Error())
		return
	}

	plan.ID = types.StringValue("placeholder")
	plan.CertificateTxtName = types.StringValue(name)
	plan.CertificateTxtValue = types.StringValue(value)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *domainSslRecordsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state domainSslRecordsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID := state.OrganizationID.ValueString()
	subdomain := state.StatusPageSubdomain.ValueString()
	statusPage, err := r.client.GetStatusPage(&orgID, &subdomain)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading status page for domain SSL records",
			"Could not read status page "+subdomain+": "+err.Error(),
		)
		return
	}

	if statusPage.DomainConfig == nil || statusPage.DomainConfig.Status == nil ||
		*statusPage.DomainConfig.Status == "disabled" {
		resp.State.RemoveResource(ctx)
		return
	}

	// The certificate_txt_* fields are only present in validation_records while
	// the SSL challenge is pending. Once the cert is issued, the API stops
	// returning them (or returns them empty). The TXT record they describe is
	// still required at the DNS provider for cert renewal, so preserve the
	// existing state values rather than clearing them — clearing would force a
	// destructive replan of the downstream cloudflare_record.txt.
	if vr := statusPage.DomainConfig.ValidationRecords; vr != nil {
		if v, ok := vr["certificate_txt_name"]; ok && v != "" {
			state.CertificateTxtName = types.StringValue(v)
		}
		if v, ok := vr["certificate_txt_value"]; ok && v != "" {
			state.CertificateTxtValue = types.StringValue(v)
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *domainSslRecordsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan domainSslRecordsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name, value, err := r.pollUntilCertRecordsReady(
		ctx,
		plan.OrganizationID.ValueString(),
		plan.StatusPageSubdomain.ValueString(),
		plan.TimeoutSeconds.ValueInt64(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Failed waiting for SSL certificate records", err.Error())
		return
	}

	plan.ID = types.StringValue("placeholder")
	plan.CertificateTxtName = types.StringValue(name)
	plan.CertificateTxtValue = types.StringValue(value)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *domainSslRecordsResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

func (r *domainSslRecordsResource) Configure(
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

// pollUntilCertRecordsReady polls GetStatusPage until certificate_txt_name appears in validation_records.
func (r *domainSslRecordsResource) pollUntilCertRecordsReady(ctx context.Context, orgID, subdomain string, timeoutSeconds int64) (name, value string, err error) {
	deadline := time.Now().Add(time.Duration(timeoutSeconds) * time.Second)

	for {
		if time.Now().After(deadline) {
			return "", "", fmt.Errorf(
				"timed out after %ds waiting for SSL certificate records on status page %q",
				timeoutSeconds, subdomain,
			)
		}

		select {
		case <-ctx.Done():
			return "", "", fmt.Errorf("context cancelled while waiting for SSL certificate records on status page %q", subdomain)
		default:
		}

		statusPage, pollErr := r.client.GetStatusPage(&orgID, &subdomain)
		if pollErr != nil {
			return "", "", fmt.Errorf("error polling status page %q: %w", subdomain, pollErr)
		}

		if statusPage.DomainConfig == nil {
			return "", "", fmt.Errorf("status page %q has no domain_config; ensure domain_config is set before using this resource", subdomain)
		}

		if vr := statusPage.DomainConfig.ValidationRecords; vr != nil {
			if n, ok := vr["certificate_txt_name"]; ok && n != "" {
				v := vr["certificate_txt_value"]
				return n, v, nil
			}
		}

		// Records not yet present — wait and retry.
		select {
		case <-ctx.Done():
			return "", "", fmt.Errorf("context cancelled while waiting for SSL certificate records on status page %q", subdomain)
		case <-time.After(sslRecordsPollInterval):
		}
	}
}
