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
