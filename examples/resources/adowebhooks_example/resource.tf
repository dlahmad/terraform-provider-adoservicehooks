resource "adowebhooks_repository_webhook" "example" {
  project_id    = "10e4f93c-42a6-4d81-9c1f-6a177953f216" // ID, not the name of the project
  repository_id = "afdf4653-27ef-4d8f-b253-15c93cc3f931" // ID, not the name of the repository
  url           = "https://example.com/webhook2"
  event_type    = "git.push"
}
