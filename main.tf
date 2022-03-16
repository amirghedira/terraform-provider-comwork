

terraform {
  required_providers {
    myprovider = {
      source  = "terraform.local/local/myprovider"
      version = "1.0.0"
    }
  }
}



provider "myprovider" {
    region = "fr"
    token = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpZCI6MSwiZW1haWwiOiJhbWlyZ2hlZGlyYTA2QGdtYWlsLmNvbSIsInRpbWUiOiIwMy8xNi8yMDIyLCAyMzozNDoxNiJ9.FiGhNPVPFPtSR7xCDFCbdcaDrELshiXR6_paW5VqqRE"
    ngx_username = "sre-devops"
    ngx_password = "QaMj8veb6RLkwPgkwb3SXBNs"
}

resource "myprovider_instance" "my_project" {
  project_name = "the_terraform_projeeect_test"
  stack_name = "the_terraform_projeceeet_test"
  project_type = "code"
  instance_type = "DEV1-S"
  email = "amirghedira06@gmail.com"

}
