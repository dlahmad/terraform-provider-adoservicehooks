// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &SubscriptionResource{}

func NewSubscriptionResource() resource.Resource {
	return &SubscriptionResource{}
}

// SubscriptionResource defines the resource implementation for a webhook in Azure DevOps.
type SubscriptionResource struct {
	client *Client
}

// Metadata returns the resource type name.
func (r *SubscriptionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subscription"
}

// Schema defines the resource schema.
func (r *SubscriptionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"consumer_action_id": schema.StringAttribute{
				Required:    true,
				Description: "The action the consumer will perform, typically representing the type of request, such as an HTTP request.",
			},
			"consumer_id": schema.StringAttribute{
				Required:    true,
				Description: "Identifies the consumer of the webhook. For example, 'webHooks' to indicate that a webhook will be triggered.",
			},
			"consumer_inputs": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Inputs that are required by the consumer action, such as URL, authentication, and headers.",
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						Optional:    true,
						Description: "The target URL for the webhook where the HTTP request will be sent.",
					},
					"basic_auth_username": schema.StringAttribute{
						Optional:    true,
						Description: "The username for basic HTTP authentication when invoking the webhook.",
					},
					"basic_auth_password": schema.StringAttribute{
						Optional:    true,
						Sensitive:   true,
						Description: "The password for basic HTTP authentication when invoking the webhook. Marked as sensitive to prevent exposure in logs.",
					},
					"http_headers": schema.StringAttribute{
						Optional:    true,
						Description: "A list of HTTP headers to include in the webhook request, formatted as a comma-separated string (e.g., 'Header1:Value1,Header2:Value2').",
					},
					"resource_details_to_send": schema.StringAttribute{
						Optional:    true,
						Description: "Specifies the level of resource detail that will be sent to the webhook (e.g., 'minimal' or 'detailed').",
					},
					"messages_to_send": schema.StringAttribute{
						Optional:    true,
						Description: "Defines which messages, if any, will be sent to the webhook. Typically 'none' to send no messages.",
					},
					"detailed_messages_to_send": schema.StringAttribute{
						Optional:    true,
						Description: "Defines whether detailed messages should be sent to the webhook, usually 'none'.",
					},
				},
			},
			"event_type": schema.StringAttribute{
				Required:    true,
				Description: "The type of event that triggers the webhook, such as 'git.push' for a Git push event.",
			},
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the webhook subscription. This is usually computed by the system.",
			},
			"publisher_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the publisher that initiates the event (e.g., 'tfs' for Azure DevOps or Team Foundation Server).",
			},
			"publisher_inputs": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Details about the publisher and the specific resources related to the event.",
				Attributes: map[string]schema.Attribute{
					"repository": schema.StringAttribute{
						Optional:    true,
						Description: "The repository from which the event (such as a push) originates.",
					},
					"branch": schema.StringAttribute{
						Optional:    true,
						Description: "The branch in the repository where the event occurred.",
					},
					"pushed_by": schema.StringAttribute{
						Optional:    true,
						Description: "The user who pushed the changes in a Git push event.",
					},
					"project_id": schema.StringAttribute{
						Optional:    true,
						Description: "The unique ID of the project associated with the event.",
					},
					"tfs_subscription_id": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The subscription ID from TFS or Azure DevOps that identifies this specific webhook subscription.",
					},
				},
			},
			"resource_version": schema.StringAttribute{
				Optional:    true,
				Description: "The version of the resource triggering the webhook event, typically set to '1.0' or another version string.",
			},
			"scope": schema.Int64Attribute{
				Optional:    true,
				Description: "Defines the scope of the webhook event. This is often an integer representing a specific scope or context.",
			},
		},
	}
}

// Configure sets up the client for the resource.
func (r *SubscriptionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *SubscriptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WebhookSubscriptionTF

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestData := ConvertToJSONModel(&data)

	// Create the webhook using the client and pass the project_id
	webhookResponse, err := r.client.CreateOrUpdateWebhook(
		requestData,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to create webhook: %s", err),
		)
		return
	}

	// Set the webhook ID after creation
	data = *ConvertToTFModel(webhookResponse)

	// Log creation
	tflog.Trace(ctx, "Created Azure DevOps Webhook")

	// Save the data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the current state of the Azure DevOps webhook resource.
func (r *SubscriptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data WebhookSubscriptionTF

	// Read prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the webhook using the client and pass the project_id
	webhookResponse, err := r.client.GetWebhook(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to get webhook: %s", err),
		)
		return
	}

	// Update the model with the current state of the webhook
	data = *ConvertToTFModel(webhookResponse)

	// Save the updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the state of an existing Azure DevOps webhook.
func (r *SubscriptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData WebhookSubscriptionTF
	var stateData WebhookSubscriptionTF

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
	planData.ID = stateData.ID
	// Log the webhook ID to verify it's being retrieved from the state correctly
	tflog.Info(ctx, "Webhook ID from state: "+stateData.ID.ValueString())

	requestData := ConvertToJSONModel(&planData)

	// Use the planData values for the updated webhook details
	webhookResponse, err := r.client.CreateOrUpdateWebhook(requestData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to update webhook: %s", err),
		)
		return
	}

	// Set the updated values (from the response) to planData
	updatedData := ConvertToTFModel(webhookResponse)

	// Save the updated data into Terraform state (from planData which now holds updated values)
	resp.Diagnostics.Append(resp.State.Set(ctx, updatedData)...)
}

// Delete deletes an existing Azure DevOps webhook.
func (r *SubscriptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data WebhookSubscriptionTF

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the webhook using the client and pass the project_id
	err := r.client.DeleteWebhook(data.ID.ValueString())
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

func (r *SubscriptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
	var data = *ConvertToTFModel(webhookResponse)
	// Set the imported state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
