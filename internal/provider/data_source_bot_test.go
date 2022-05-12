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
						"data.awslex_bot_resource.foo", "id", regexp.MustCompile("C5H22UIPWC")),
				),
			},
		},
	})
}

// use id and alias from the stable dev bot
const testAccDataSourceBot = `
provider "awslex" {
  region = "us-west-2"
}

data "awslex_bot_resource" "foo" {
  id = "C5H22UIPWC"
  alias = "latest"
}
`
