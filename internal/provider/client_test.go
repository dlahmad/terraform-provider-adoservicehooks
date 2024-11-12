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

	resN, err := client.CreateOrUpdateWebhook(project, repository, hookUrl, "git.push", nil)
	if err != nil || resN.ID == "" || resN.EventType == "" {
		t.Fail()
	}

	resN, err = client.GetWebhook(resN.ID)
	if err != nil || resN.ID == "" || resN.EventType == "" {
		t.Fail()
	}

	resN, err = client.CreateOrUpdateWebhook(project, repository, hookUrl, "git.push", &resN.ID)
	if err != nil || resN.ID == "" {
		t.Fail()
	}

	err = client.DeleteWebhook(project, resN.ID)
	if err != nil {
		t.Fail()
	}
	// Your test code here...
}
