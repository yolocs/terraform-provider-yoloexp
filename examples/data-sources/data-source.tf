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

data "yoloexp_notion_database" "example" {
  id = "53172ac8-a84f-4b06-b3c4-81f501f2d46c"
}

output "properties" {
  value = data.yoloexp_notion_database.example.properties
}
