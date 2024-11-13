// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	org := os.Getenv("ADO_ORGANIZATION")
	pat := os.Getenv("ADO_PAT")
	project := os.Getenv("ADO_PROJECT")
	repository := os.Getenv("ADO_REPOSITORY")
	hookUrl := os.Getenv("ADO_HOOK_URL")

	if org == "" || pat == "" || project == "" || repository == "" || hookUrl == "" {
		return
	}

	// Continue with the test, e.g., setting up an Azure DevOps client
	// Example: Create client and test against Azure DevOps API
	client, err := NewClient(&org, &pat)

	if err != nil {
		t.Fail()
	}

	res, err := client.GetProjectGuid(project)

	if err != nil || res.ID == "" {
		t.Fail()
	}

	res, err = client.GetRepositoryGuid(project, repository)

	if err != nil || res.ID == "" {
		t.Fail()
	}

	subscription := DefaultWebhookSubscription()
	subscription.PublisherInputs = &PublisherInputs{
		RepositoryId: &repository,
		Branch:       stringToPointer("master"),
		ProjectId:    &project,
	}
	subscription.ConsumerInputs = &ConsumerInputs{
		URL: stringToPointer(hookUrl),
	}

	resN, err := client.CreateOrUpdateWebhook(subscription)
	if err != nil || resN.ID == nil || resN.EventType == nil {
		t.Fail()
	}

	resN, err = client.GetWebhook(*resN.ID)
	if err != nil || resN.ID == nil || resN.EventType == nil {
		t.Fail()
	}

	resN, err = client.CreateOrUpdateWebhook(resN)
	if err != nil || resN.ID == nil {
		t.Fail()
	}

	err = client.DeleteWebhook(*resN.ID)
	if err != nil {
		t.Fail()
	}
	// Your test code here...
}
