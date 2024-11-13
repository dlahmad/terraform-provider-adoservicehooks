// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type ConsumerInputsTF struct {
	URL                    types.String `tfsdk:"url"`
	BasicAuthUsername      types.String `tfsdk:"basic_auth_username"`
	BasicAuthPassword      types.String `tfsdk:"basic_auth_password"`
	HTTPHeaders            types.String `tfsdk:"http_headers"`
	ResourceDetailsToSend  types.String `tfsdk:"resource_details_to_send"`
	MessagesToSend         types.String `tfsdk:"messages_to_send"`
	DetailedMessagesToSend types.String `tfsdk:"detailed_messages_to_send"`
}

type PublisherInputsTF struct {
	RepositoryId      types.String `tfsdk:"repository"`
	Branch            types.String `tfsdk:"branch"`
	PushedBy          types.String `tfsdk:"pushed_by"`
	ProjectId         types.String `tfsdk:"project_id"`
	TfsSubscriptionId types.String `tfsdk:"tfs_subscription_id"`
}

type WebhookSubscriptionTF struct {
	ConsumerActionId types.String       `tfsdk:"consumer_action_id"`
	ConsumerId       types.String       `tfsdk:"consumer_id"`
	ConsumerInputs   *ConsumerInputsTF  `tfsdk:"consumer_inputs"`
	EventType        types.String       `tfsdk:"event_type"`
	ID               types.String       `tfsdk:"id"`
	PublisherId      types.String       `tfsdk:"publisher_id"`
	PublisherInputs  *PublisherInputsTF `tfsdk:"publisher_inputs"`
	ResourceVersion  types.String       `tfsdk:"resource_version"`
	Scope            types.Int64        `tfsdk:"scope"`
}

func stringToPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func int64Pointer(i int64) *int64 {
	return &i
}

func DefaultWebhookSubscriptionTF() *WebhookSubscriptionTF {
	return &WebhookSubscriptionTF{
		ConsumerId:      types.StringValue("webHooks"),
		ResourceVersion: types.StringValue("1.0"),
		Scope:           types.Int64Value(1),
	}
}

func (ws *WebhookSubscriptionTF) SetDefaults() {
	if ws.ConsumerId.IsNull() || ws.ConsumerId.IsUnknown() {
		ws.ConsumerId = types.StringValue("webHooks")
	}
	if ws.ResourceVersion.IsNull() || ws.ResourceVersion.IsUnknown() {
		ws.ResourceVersion = types.StringValue("1.0")
	}
	if ws.Scope.IsNull() || ws.Scope.IsUnknown() {
		ws.Scope = types.Int64Value(1)
	}
}

func (ws *WebhookSubscription) SetDefaults() {
	// Set default value for ConsumerId if not set
	if ws.ConsumerId == "" {
		ws.ConsumerId = "webHooks"
	}

	// Set default value for PublisherId if not set
	if ws.PublisherId == nil {
		ws.PublisherId = stringToPointer("tfs")
	}

	// Set default value for ResourceVersion if not set
	if ws.ResourceVersion == nil {
		ws.ResourceVersion = stringToPointer("1.0")
	}

	// Set default value for Scope if not set
	if ws.Scope == nil {
		ws.Scope = int64Pointer(1)
	}
}

type ConsumerInputs struct {
	URL                    *string `json:"url,omitempty"`
	BasicAuthUsername      *string `json:"basicAuthUsername,omitempty"`
	BasicAuthPassword      *string `json:"basicAuthPassword,omitempty"`
	HTTPHeaders            *string `json:"httpHeaders,omitempty"`
	ResourceDetailsToSend  *string `json:"resourceDetailsToSend,omitempty"`
	MessagesToSend         *string `json:"messagesToSend,omitempty"`
	DetailedMessagesToSend *string `json:"detailedMessagesToSend,omitempty"`
}

type PublisherInputs struct {
	RepositoryId      *string `json:"repository,omitempty"`
	Branch            *string `json:"branch,omitempty"`
	PushedBy          *string `json:"pushedBy,omitempty"`
	ProjectId         *string `json:"projectId,omitempty"`
	TfsSubscriptionId *string `json:"tfsSubscriptionId,omitempty"`
}

type WebhookSubscription struct {
	ConsumerActionId *string          `json:"consumerActionId"`
	ConsumerId       string           `json:"consumerId"`
	ConsumerInputs   *ConsumerInputs  `json:"consumerInputs,omitempty"`
	EventType        *string          `json:"eventType"`
	ID               *string          `json:"id,omitempty"`
	PublisherId      *string          `json:"publisherId,omitempty"`
	PublisherInputs  *PublisherInputs `json:"publisherInputs,omitempty"`
	ResourceVersion  *string          `json:"resourceVersion,omitempty"`
	Scope            *int64           `json:"scope,omitempty"`
}

func DefaultWebhookSubscription() *WebhookSubscription {
	return &WebhookSubscription{
		ConsumerId:      "webHooks",
		PublisherId:     stringToPointer("tfs"),
		ResourceVersion: stringToPointer("1.0"),
		Scope:           int64Pointer(1),
	}
}

func ConvertToJSONModel(tf *WebhookSubscriptionTF) *WebhookSubscription {
	return &WebhookSubscription{
		ConsumerActionId: getOptionalString(tf.ConsumerActionId),
		ConsumerId:       tf.ConsumerId.ValueString(),
		ConsumerInputs: &ConsumerInputs{
			URL:                    getOptionalString(tf.ConsumerInputs.URL),
			BasicAuthUsername:      getOptionalString(tf.ConsumerInputs.BasicAuthUsername),
			BasicAuthPassword:      getOptionalString(tf.ConsumerInputs.BasicAuthPassword),
			HTTPHeaders:            getOptionalString(tf.ConsumerInputs.HTTPHeaders),
			ResourceDetailsToSend:  getOptionalString(tf.ConsumerInputs.ResourceDetailsToSend),
			MessagesToSend:         getOptionalString(tf.ConsumerInputs.MessagesToSend),
			DetailedMessagesToSend: getOptionalString(tf.ConsumerInputs.DetailedMessagesToSend),
		},
		EventType:   getOptionalString(tf.EventType),
		ID:          getOptionalString(tf.ID),
		PublisherId: getOptionalString(tf.PublisherId),
		PublisherInputs: &PublisherInputs{
			RepositoryId:      getOptionalString(tf.PublisherInputs.RepositoryId),
			Branch:            getOptionalString(tf.PublisherInputs.Branch),
			PushedBy:          getOptionalString(tf.PublisherInputs.PushedBy),
			ProjectId:         getOptionalString(tf.PublisherInputs.ProjectId),
			TfsSubscriptionId: getOptionalString(tf.PublisherInputs.TfsSubscriptionId),
		},
		ResourceVersion: getOptionalString(tf.ResourceVersion),
		Scope:           getOptionalInt64(tf.Scope),
	}
}

// Convert from JSON structs to Terraform SDK structs
func ConvertToTFModel(json *WebhookSubscription) *WebhookSubscriptionTF {
	return &WebhookSubscriptionTF{
		ConsumerActionId: types.StringPointerValue(json.ConsumerActionId),
		ConsumerId:       types.StringValue(json.ConsumerId),
		ConsumerInputs: &ConsumerInputsTF{
			URL:                    types.StringPointerValue(json.ConsumerInputs.URL),
			BasicAuthUsername:      types.StringPointerValue(json.ConsumerInputs.BasicAuthUsername),
			BasicAuthPassword:      types.StringPointerValue(json.ConsumerInputs.BasicAuthPassword),
			HTTPHeaders:            types.StringPointerValue(json.ConsumerInputs.HTTPHeaders),
			ResourceDetailsToSend:  types.StringPointerValue(json.ConsumerInputs.ResourceDetailsToSend),
			MessagesToSend:         types.StringPointerValue(json.ConsumerInputs.MessagesToSend),
			DetailedMessagesToSend: types.StringPointerValue(json.ConsumerInputs.DetailedMessagesToSend),
		},
		EventType:   types.StringPointerValue(json.EventType),
		ID:          types.StringPointerValue(json.ID),
		PublisherId: types.StringPointerValue(json.PublisherId),
		PublisherInputs: &PublisherInputsTF{
			RepositoryId:      types.StringPointerValue(json.PublisherInputs.RepositoryId),
			Branch:            types.StringPointerValue(json.PublisherInputs.Branch),
			PushedBy:          types.StringPointerValue(json.PublisherInputs.PushedBy),
			ProjectId:         types.StringPointerValue(json.PublisherInputs.ProjectId),
			TfsSubscriptionId: types.StringPointerValue(json.PublisherInputs.TfsSubscriptionId),
		},
		ResourceVersion: types.StringPointerValue(json.ResourceVersion),
		Scope:           types.Int64PointerValue(json.Scope),
	}
}

// Helper function to handle optional strings in Terraform SDK
func getOptionalString(t types.String) *string {
	if t.IsNull() || t.IsUnknown() {
		return nil
	}
	val := t.ValueString()
	return &val
}

// Helper function to handle optional int64 in Terraform SDK
func getOptionalInt64(t types.Int64) *int64 {
	if t.IsNull() || t.IsUnknown() {
		return nil
	}
	val := t.ValueInt64()
	return &val
}
