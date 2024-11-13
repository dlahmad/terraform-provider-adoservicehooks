// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	HTTPClient   *http.Client
	Organization string
	Pat          string
	BaseURL      string
}

func NewClient(organization, pat *string) (*Client, error) {
	c := &Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		// Assuming a default URL for Azure DevOps organization
		BaseURL:      "https://dev.azure.com/",
		Organization: "",
		Pat:          "",
	}

	if organization != nil {
		c.Organization = *organization
	}

	if pat != nil {
		c.Pat = *pat
	}

	// If any required field is not provided, return empty client
	if organization == nil || pat == nil {
		return c, nil
	}

	return c, nil
}

func (c *Client) createRawRequest(method, url string, body interface{}) (*http.Request, error) {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	// Use project as a parameter for the URL
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Adding Authorization header using the Personal Access Token (PAT)
	req.SetBasicAuth("", c.Pat)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (c *Client) GetProjectGuid(project string) (*IdResponse, error) {
	req, err := c.createRawRequest("GET", c.BaseURL+c.Organization+"/_apis/projects/"+project+"?api-version=7.0", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get id, status code: %d", resp.StatusCode)
	}

	var webhookResponse IdResponse
	if err := json.NewDecoder(resp.Body).Decode(&webhookResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &webhookResponse, nil
}

func (c *Client) GetRepositoryGuid(project, repository string) (*IdResponse, error) {
	req, err := c.createRawRequest("GET", c.BaseURL+c.Organization+"/"+project+"/_apis/git/repositories/"+repository+"?api-version=7.0", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get id, status code: %d", resp.StatusCode)
	}

	var webhookResponse IdResponse
	if err := json.NewDecoder(resp.Body).Decode(&webhookResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &webhookResponse, nil
}

func (c *Client) GetWebhook(webhookID string) (*WebhookSubscription, error) {
	req, err := c.createRawRequest("GET", c.BaseURL+c.Organization+"/_apis/hooks/subscriptions/"+webhookID+"/?api-version=7.0", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get webhook, status code: %d", resp.StatusCode)
	}

	var webhookResponse WebhookSubscription
	if err := json.NewDecoder(resp.Body).Decode(&webhookResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &webhookResponse, nil
}

func (c *Client) CreateOrUpdateWebhook(subscription *WebhookSubscription) (*WebhookSubscription, error) {

	verb := "POST"

	// Simulating a ternary operation
	if subscription.ID != nil {
		verb = "PUT"
	}

	// Create the request
	req, err := c.createRawRequest(verb, c.BaseURL+c.Organization+"/_apis/hooks/subscriptions?api-version=7.0", subscription)
	if err != nil {
		return nil, err
	}

	// Send the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the status code is created
	if resp.StatusCode < 200 || resp.StatusCode > 201 {
		return nil, fmt.Errorf("failed to create webhook, status code: %d", resp.StatusCode)
	}

	// Parse the response into WebhookResponse
	var webhookResponse WebhookSubscription
	if err := json.NewDecoder(resp.Body).Decode(&webhookResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &webhookResponse, nil
}

func (c *Client) DeleteWebhook(webhookID string) error {
	req, err := c.createRawRequest("DELETE", c.BaseURL+c.Organization+"/_apis/hooks/subscriptions/"+webhookID+"?api-version=7.0", nil)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete webhook, status code: %d", resp.StatusCode)
	}

	return nil
}

type IdResponse struct {
	ID string `json:"id"`
}
