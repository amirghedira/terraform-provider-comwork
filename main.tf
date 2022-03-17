

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
    token = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpZCI6MSwiZW1haWwiOiJhbWlyZ2hlZGlyYTA2QGdtYWlsLmNvbSIsInRpbWUiOiIwMy8xNy8yMDIyLCAwOToyMDowMiJ9.p1XRriQJhUSrXiPIfXg-cO3u6s_sKqZ3_7yRwNXVDXk"
    ngx_username = "sre-devops"
    ngx_password = "QaMj8veb6RLkwPgkwb3SXBNs"
}

resource "comwork_instance" "my_project" {
  project_name = "terraform_testing_provision"
  stack_name = "terraform_testing_provision"
  project_type = "code"
  status = "poweroff"
  instance_type = "DEV1-S"
  email = "amirghedira06@gmail.com"

}
