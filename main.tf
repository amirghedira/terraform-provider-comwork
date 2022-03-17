

terraform {
  required_providers {
    comwork = {
      source  = "terraform.local/local/comwork"
      version = "1.0.0"
    }
  }
}



provider "comwork" {
    region = "fr-par-1"
    token = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpZCI6MSwiZW1haWwiOiJhbWlyZ2hlZGlyYTA2QGdtYWlsLmNvbSIsInRpbWUiOiIwMy8xNy8yMDIyLCAwMDoyMzoyMCJ9.9U8gGaLahr8ZwYlaU0fzPDf-P4ynXT-zfC-5nS-pCeA"
    ngx_username = "sre-devops"
    ngx_password = "QaMj8veb6RLkwPgkwb3SXBNs"
}

resource "comwork_instance" "my_project" {
  project_name = "terraform_action"
  stack_name = "terraform_action"
  project_type = "code"
  status = "poweron"
  instance_type = "DEV1-S"
  email = "amirghedira06@gmail.com"

}
