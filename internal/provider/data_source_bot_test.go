package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceBot(t *testing.T) {
	// t.Skip("data source not yet implemented, remove this once you add your own code")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBot,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.awslex_bot_resource.foo", "id", regexp.MustCompile("QU1ORIZZTP")),
				),
			},
		},
	})
}

const testAccDataSourceBot = `
provider "awslex" {
  region = "us-west-2"
}

data "awslex_bot_resource" "foo" {
  id = "QU1ORIZZTP"
  alias = "alias1"
}
`
