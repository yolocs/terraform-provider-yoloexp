// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestNotionPageDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNotionPageDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.yoloexp_notion_page.test", "id", "example-id"),
				),
			},
		},
	})
}

func TestNotionDatabaseDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNotionDatabaseDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.yoloexp_notion_database.test", "id", "example-id"),
				),
			},
		},
	})
}

const (
	testAccNotionPageDataSourceConfig = `
data "yoloexp_notion_page" "test" {
  id = "example-id"
}
`

	testAccNotionDatabaseDataSourceConfig = `
data "yoloexp_notion_page" "test" {
  id = "example-id"
}
`
)
