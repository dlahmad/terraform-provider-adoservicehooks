# adoservicehooks Provider ![Build Status](https://github.com/dlahmad/terraform-provider-adoservicehooks/actions/workflows/release.yml/badge.svg)

A little provider allowing to register webhooks (for now only based on repository based triggers) for AzureDevOps. Unfortunately the regular [AzureDevOps Terraform Provider](https://registry.terraform.io/providers/microsoft/azuredevops/latest/docs) does not support registering service hooks. For now, this provider supports the state managed usage of [Service Hook Subscription API](https://learn.microsoft.com/en-us/rest/api/azure/devops/hooks/subscriptions/create?view=azure-devops-rest-7.1&tabs=HTTP) in a flexibel (pass through with no additional validation) manner.

## Example Usage

```terraform
provider "adoservicehooks" {
  organization = "yourorg"
  pat          = "<pat>"
}

resource "adoservicehooks_subscription" "example" {
  consumer_action_id = "httpRequest"
  consumer_id        = "webHooks"
  consumer_inputs = {
    url          = "https://triggerservice.com/webhook"
    http_headers = "TRIGGERSOURCE:DEVOPS"
  } project
  event_type   = "git.push"
  publisher_id = "tfs"
  publisher_inputs = {
    project_id = "3d994ea0-b3c3-4fca-8318-91ec3d042b9d" 
    repository = "760d53d4-f394-44fc-ab00-5932b6b7da9d" 
  }
}
```

To find out which options you have check the [adoservicehooks Provider Documentation](https://registry.terraform.io/providers/dlahmad/adoservicehooks/latest/docs) and [Service Hook Subscription API](https://learn.microsoft.com/en-us/rest/api/azure/devops/hooks/subscriptions/create?view=azure-devops-rest-7.1&tabs=HTTP).
