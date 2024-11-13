// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &azureDevopsWebhooksProvider{}
)

type azureDevopsWebhooksProviderModel struct {
	Organization types.String `tfsdk:"organization"`
	Pat          types.String `tfsdk:"pat"`
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &azureDevopsWebhooksProvider{
			version: version,
		}
	}
}

// azureDevopsWebhooksProvider is the provider implementation.
type azureDevopsWebhooksProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *azureDevopsWebhooksProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "adowebhooks"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *azureDevopsWebhooksProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				Optional: true,
			},
			"pat": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// Configure prepares a HashiCups API client for data sources and resources.
func (p *azureDevopsWebhooksProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config azureDevopsWebhooksProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Organization.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("organization"),
			"Unknown AzureDevOps Organization",
			"The provider cannot create the client because it needs to know the AzureDevOps organization. "+
				"Set the organization value in the configuration or use the ADOWEBHOOKS_ORGANIZATION environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if config.Pat.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("pat"),
			"Unknown AzureDevOps PAT",
			"The provider cannot create the client because it needs to know the AzureDevOps PAT. "+
				"Set the password value in the configuration or use the ADOWEBHOOKS_PAT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	organization := os.Getenv("ADOWEBHOOKS_ORGANIZATION")
	pat := os.Getenv("ADOWEBHOOKS_PAT")

	if !config.Organization.IsNull() {
		organization = config.Organization.ValueString()
	}

	if !config.Pat.IsNull() {
		pat = config.Pat.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if organization == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("organization"),
			"Missing AzureDevOps Organization",
			"The provider cannot create the client because it needs to know the AzureDevOps organization. "+
				"Set the organization value in the configuration or use the ADOWEBHOOKS_ORGANIZATION environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if pat == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("pat"),
			"Missing AzureDevOps PAT",
			"The provider cannot create the client because it needs to know the AzureDevOps PAT. "+
				"Set the password value in the configuration or use the ADOWEBHOOKS_PAT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new HashiCups client using the configuration values
	client, err := NewClient(&organization, &pat)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Client",
			"An unexpected error occurred when creating the client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Client Error: "+err.Error(),
		)
		return
	}

	// Make the HashiCups client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *azureDevopsWebhooksProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources defines the resources implemented in the provider.
func (p *azureDevopsWebhooksProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSubscriptionResource,
	}
}
