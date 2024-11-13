resource "adoservicehooks_subscription" "example" {
  consumer_action_id = "httpRequest"
  consumer_id        = "webHooks"
  consumer_inputs = {
    url = "https://example.com/webhook"
  }
  event_type   = "git.push"
  publisher_id = "tfs"
  publisher_inputs = {
    project_id = "10e4f93c-42a6-4d81-9c1f-6a177953f216"
    repository = "afdf4653-27ef-4d8f-b253-15c93cc3f931"
  }
}
