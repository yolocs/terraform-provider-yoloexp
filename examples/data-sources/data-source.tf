# Needs to update ~/.terraformrc to support local development.
terraform {
  required_providers {
    yoloexp = {
      source = "registry.terraform.io/hashicorp/yoloexp"
    }
  }
}

provider "yoloexp" {
  # example configuration here
}

data "yoloexp_notion_page" "example" {
  # Replace with your own page id.
  id = "8263e830-3424-475a-801c-1d971606cd6c"
}
