package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSkills(t *testing.T) {
	t.Skip("TODO: test not yet working")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSkills,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"awslex_bot_resource.foo", "id", regexp.MustCompile("TBD.*")),
				),
			},
		},
	})
}

const testAccResourceSkills = `
provider "awslex" {
  region = "us-west-2"
}

resource "awslex_bot_resource" "socal_gas_qnabot" {

  name       = "integration-bot"

  description = "bot created by integration tests"

  # path to the bot sources zip file, in bot import/export format
	# TODO: determine what path to use here
  archive_path = "TBD"

  # detect changes to the bot sources and update the bot
  source_code_hash = "TBD"

  # version of the bot
  alias = "latest"

  # arn of the lambda that fulfills the bot intents
  lambda_arn = "arn:aws:lambda:us-west-2:111365482541:function:scg-geeou-dev-wus2-lambda-fulfillment-dev"

  iam_role = "arn:aws:iam::111365482541:role/scg-lexbot-dev-wus2-iam-role-qnabot-dev"

  tags = {
    name                = "scg-shcva Virtual Assistant"
    tag-version         = "1.0.0"
    unit                = "shcva"
    portfolio           = "geeou"
    support-group       = "SCGMA Team"
    cmdb-ci-id          = "7777777"
    data-classification = "testing"
  }
}
`
