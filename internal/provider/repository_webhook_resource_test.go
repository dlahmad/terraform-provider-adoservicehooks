// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRepositoryWebhookResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },   // PreCheck to validate environment
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories, // Setting up the provider factories
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccRepositoryWebhookResourceConfig("https://example.com/webhook", "git.push", "some-repo"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("adowebhooks_webhook.test", "url", "https://example.com/webhook"),
					resource.TestCheckResourceAttr("adowebhooks_webhook.test", "event_type", "git.push"),
					resource.TestCheckResourceAttr("adowebhooks_webhook.test", "repository_id", "some-repo"),
					// Assuming "webhook_id" will be assigned
					resource.TestCheckResourceAttrSet("adowebhooks_webhook.test", "webhook_id"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "adowebhooks_webhook.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"url", "event_type"}, // Ignore these for this test
			},
			// Update and Read testing
			{
				Config: testAccRepositoryWebhookResourceConfig("https://example.com/updated_webhook", "git.push", "updated-repo"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("adowebhooks_webhook.test", "url", "https://example.com/updated_webhook"),
					resource.TestCheckResourceAttr("adowebhooks_webhook.test", "repository_id", "updated-repo"),
				),
			},
		},
	})
}

// Helper function to generate the resource configuration for testing
func testAccRepositoryWebhookResourceConfig(url, eventType, repositoryId string) string {
	return fmt.Sprintf(`
resource "adowebhooks_webhook" "test" {
  project_id    = "test-project"
  url           = "%s"
  repository_id = "%s"
  event_type    = "%s"
}
`, url, repositoryId, eventType)
}
