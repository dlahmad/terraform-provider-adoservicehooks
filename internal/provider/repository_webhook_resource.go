// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &RepositoryWebhookResource{}

func NewRepositoryWebhookResource() resource.Resource {
	return &RepositoryWebhookResource{}
}

// RepositoryWebhookResource defines the resource implementation for a webhook in Azure DevOps.
type RepositoryWebhookResource struct {
	client *Client
}

// RepositoryWebhookModel describes the resource data model.
type RepositoryWebhookModel struct {
	ProjectId    types.String `tfsdk:"project_id"` // Added project_id as a required field
	WebhookId    types.String `tfsdk:"webhook_id"`
	URL          types.String `tfsdk:"url"`
	RepositoryId types.String `tfsdk:"repository_id"`
	EventType    types.String `tfsdk:"event_type"`
}

// Metadata returns the resource type name.
func (r *RepositoryWebhookResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_webhook"
}

// Schema defines the resource schema.
func (r *RepositoryWebhookResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A resource representing a webhook in Azure DevOps.",

		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The Azure DevOps project name.",
				Required:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The URL that receives webhook notifications.",
				Required:            true,
			},
			"repository_id": schema.StringAttribute{
				MarkdownDescription: "The name of the repository to which the webhook is attached.",
				Required:            true,
			},
			"event_type": schema.StringAttribute{
				MarkdownDescription: "The event type that triggers the webhook.",
				Required:            true,
			},
			"webhook_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the webhook.",
				Computed:            true,
			},
		},
	}
}

// Configure sets up the client for the resource.
func (r *RepositoryWebhookResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	// The client is passed through ConfigureRequest.ResourceData
	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Client Configuration Error",
			"Expected *Client, got something else",
		)
		return
	}

	r.client = client
}

// Create creates a new Azure DevOps webhook using the provided parameters.
func (r *RepositoryWebhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RepositoryWebhookModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the webhook using the client and pass the project_id
	webhookResponse, err := r.client.CreateOrUpdateWebhook(
		data.ProjectId.ValueString(),
		data.RepositoryId.ValueString(),
		data.URL.ValueString(),
		data.EventType.ValueString(),
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to create webhook: %s", err),
		)
		return
	}

	// Set the webhook ID after creation
	data.WebhookId = types.StringValue(webhookResponse.ID)
	data.URL = types.StringValue(webhookResponse.ConsumerInputs.URL)
	data.EventType = types.StringValue(webhookResponse.EventType)
	data.RepositoryId = types.StringValue(webhookResponse.PublisherInputs.Repository)
	data.ProjectId = types.StringValue(webhookResponse.PublisherInputs.ProjectId)

	// Log creation
	tflog.Trace(ctx, "Created Azure DevOps Webhook")

	// Save the data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the current state of the Azure DevOps webhook resource.
func (r *RepositoryWebhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RepositoryWebhookModel

	// Read prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the webhook using the client and pass the project_id
	webhookResponse, err := r.client.GetWebhook(data.WebhookId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to get webhook: %s", err),
		)
		return
	}

	// Update the model with the current state of the webhook
	data.URL = types.StringValue(webhookResponse.ConsumerInputs.URL)
	data.EventType = types.StringValue(webhookResponse.EventType)
	data.RepositoryId = types.StringValue(webhookResponse.PublisherInputs.Repository)
	data.ProjectId = types.StringValue(webhookResponse.PublisherInputs.ProjectId)

	// Save the updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the state of an existing Azure DevOps webhook.
func (r *RepositoryWebhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData RepositoryWebhookModel
	var stateData RepositoryWebhookModel

	// Read Terraform plan data (the desired new state) into the planData model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the current state (to get the existing webhook_id) into stateData model
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the webhook ID from the current state
	webhookID := stateData.WebhookId.ValueString()

	// Log the webhook ID to verify it's being retrieved from the state correctly
	tflog.Info(ctx, "Webhook ID from state: "+webhookID)

	// Use the planData values for the updated webhook details
	webhookResponse, err := r.client.CreateOrUpdateWebhook(
		planData.ProjectId.ValueString(),    // From plan
		planData.RepositoryId.ValueString(), // From plan
		planData.URL.ValueString(),          // From plan
		planData.EventType.ValueString(),    // From plan
		&webhookID,                          // From state
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to update webhook: %s", err),
		)
		return
	}

	// Set the updated values (from the response) to planData
	planData.WebhookId = types.StringValue(webhookResponse.ID)
	planData.URL = types.StringValue(webhookResponse.ConsumerInputs.URL)
	planData.EventType = types.StringValue(webhookResponse.EventType)
	planData.RepositoryId = types.StringValue(webhookResponse.PublisherInputs.Repository)
	planData.ProjectId = types.StringValue(webhookResponse.PublisherInputs.ProjectId)

	// Save the updated data into Terraform state (from planData which now holds updated values)
	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

// Delete deletes an existing Azure DevOps webhook.
func (r *RepositoryWebhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RepositoryWebhookModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the webhook using the client and pass the project_id
	err := r.client.DeleteWebhook(data.ProjectId.ValueString(), data.WebhookId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to delete webhook: %s", err),
		)
		return
	}

	// Log the deletion
	tflog.Trace(ctx, "Deleted Azure DevOps Webhook")
}

func (r *RepositoryWebhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve the webhook ID from the request
	webhookID := req.ID

	// Use the client to fetch the webhook details using the ID
	webhookResponse, err := r.client.GetWebhook(webhookID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to get webhook: %s", err),
		)
		return
	}

	// Map the response to the model
	var data RepositoryWebhookModel
	data.WebhookId = types.StringValue(webhookResponse.ID)
	data.URL = types.StringValue(webhookResponse.ConsumerInputs.URL)
	data.EventType = types.StringValue(webhookResponse.EventType)
	data.RepositoryId = types.StringValue(webhookResponse.PublisherInputs.Repository)
	data.ProjectId = types.StringValue(webhookResponse.PublisherInputs.ProjectId)

	// Set the imported state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
