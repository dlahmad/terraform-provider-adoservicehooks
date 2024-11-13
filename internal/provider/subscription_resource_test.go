// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

// func TestAccSubscriptionResource(t *testing.T) {
// 	// Replace with actual values for testing
// 	org := "your-organization-name"
// 	pat := "your-personal-access-token"

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { testAccPreCheck(t) },   // PreCheck to validate environment
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories, // Setting up the provider factories
// 		Steps: []resource.TestStep{
// 			// Create and Read testing
// 			{
// 				Config: testAccSubscriptionResourceConfig(
// 					"git.push",
// 					"webHooks",
// 					"https://example.com/webhook",
// 					"some-repo",
// 					"some-branch",
// 					"some-user",
// 					"test-project",
// 					"publisher-id",
// 					org,
// 					pat,
// 				),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr("adoservicehooks_subscription.test", "event_type", "git.push"),
// 					resource.TestCheckResourceAttr("adoservicehooks_subscription.test", "consumer_id", "webHooks"),
// 					resource.TestCheckResourceAttr("adoservicehooks_subscription.test", "consumer_inputs.url", "https://example.com/webhook"),
// 					resource.TestCheckResourceAttr("adoservicehooks_subscription.test", "publisher_id", "publisher-id"),
// 					resource.TestCheckResourceAttr("adoservicehooks_subscription.test", "publisher_inputs.repository", "some-repo"),
// 					resource.TestCheckResourceAttr("adoservicehooks_subscription.test", "publisher_inputs.branch", "some-branch"),
// 					resource.TestCheckResourceAttr("adoservicehooks_subscription.test", "publisher_inputs.pushed_by", "some-user"),
// 					resource.TestCheckResourceAttr("adoservicehooks_subscription.test", "publisher_inputs.project_id", "test-project"),
// 					// Assuming "id" and other fields will be set correctly
// 					resource.TestCheckResourceAttrSet("adoservicehooks_subscription.test", "id"),
// 				),
// 			},
// 			// ImportState testing
// 			{
// 				ResourceName:            "adoservicehooks_subscription.test",
// 				ImportState:             true,
// 				ImportStateVerify:       true,
// 				ImportStateVerifyIgnore: []string{"consumer_inputs", "publisher_inputs"}, // Ignore these for this test
// 			},
// 			// Update and Read testing
// 			{
// 				Config: testAccSubscriptionResourceConfig(
// 					"git.push",
// 					"webHooks",
// 					"https://example.com/updated_webhook",
// 					"updated-repo",
// 					"updated-branch",
// 					"updated-user",
// 					"updated-project",
// 					"updated-publisher-id",
// 					org,
// 					pat,
// 				),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr("adoservicehooks_subscription.test", "consumer_id", "webHooks"),
// 					resource.TestCheckResourceAttr("adoservicehooks_subscription.test", "consumer_inputs.url", "https://example.com/updated_webhook"),
// 					resource.TestCheckResourceAttr("adoservicehooks_subscription.test", "publisher_id", "updated-publisher-id"),
// 					resource.TestCheckResourceAttr("adoservicehooks_subscription.test", "publisher_inputs.repository", "updated-repo"),
// 					resource.TestCheckResourceAttr("adoservicehooks_subscription.test", "publisher_inputs.branch", "updated-branch"),
// 					resource.TestCheckResourceAttr("adoservicehooks_subscription.test", "publisher_inputs.pushed_by", "updated-user"),
// 					resource.TestCheckResourceAttr("adoservicehooks_subscription.test", "publisher_inputs.project_id", "updated-project"),
// 				),
// 			},
// 		},
// 	})
// }

// // Helper function to generate the resource configuration for testing.
// func testAccSubscriptionResourceConfig(
// 	eventType, consumerId, url, repository, branch, pushedBy, projectId, publisherId, org, pat string,
// ) string {
// 	return fmt.Sprintf(`
// provider "adoservicehooks" {
//   organization = "%s"
//   pat          = "%s"
// }

// resource "adoservicehooks_subscription" "test" {
//   consumer_action_id = "some-action"
//   consumer_id        = "%s"
//   consumer_inputs = {
//     url                    = "%s"
//     basic_auth_username    = "user"
//     basic_auth_password    = "password"
//     http_headers           = "Header1:Value1,Header2:Value2"
//     resource_details_to_send = "minimal"
//     messages_to_send         = "none"
//     detailed_messages_to_send = "none"
//   }
//   event_type         = "%s"
//   publisher_id       = "%s"
//   publisher_inputs = {
//     repository        = "%s"
//     branch            = "%s"
//     pushed_by         = "%s"
//     project_id        = "%s"
//   }
//   resource_version   = "1.0"
//   scope              = 1
// }
// `, org, pat, consumerId, url, eventType, publisherId, repository, branch, pushedBy, projectId)
// }
