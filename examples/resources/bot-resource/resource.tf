terraform {
  required_providers {
    # for testing versions existing locally
    awslex = {
      source  = "localhost/va/awslex"
      version = "0.2.0-beta3"
    }

    # for testing against a release in the tfc private registry
    # awslex = {
    #   source = "app.terraform.io/SempraUtilities/awslex"
    #   version = "0.2.0-beta3"
    # }
  }
}

provider "awslex" {
  region = "us-west-2"
  // this should error out when testing locally
  // role_arn = "arn:aws:iam::123456789012:role/awslex-role"
}

provider "aws" {
  region = "us-west-2"
}

data "aws_caller_identity" "current" {}

locals {
  bot_name = "TerraBot"
  bot_description      = "Terraform Bot"
  lambda_arn           = "arn:aws:lambda:us-west-2:111365482541:function:scg-geeou-dev-wus2-lambda-fulfillment-dev"
  account_id           = data.aws_caller_identity.current.account_id
  bot_id               = awslex_bot_resource.socal_gas_qnabot.id
  bot_alias_id         = awslex_bot_resource.socal_gas_qnabot.alias_id
}

# create the file that represents the Lex bot sources
module "bot_sources" {
  source          = "./sources"
  bot_description = local.bot_description
  bot_name = local.bot_name
  intents = [
    {
      id = "gas-leak"
      questions = ["I smell gas in my house. What should I do?",
        "help my gas is leaking",
        "smell garlic in my home and I'm not cooking",
      "emergency gas leak"]

      answer = "For Gas Emergencies or Safety Issues call Emergencies: 911 For general safety issues: 1-800-427-2200"
    },
    {
      id = "password-reset"
      questions = ["how do I reset my password",
        "I forgot my password",
        "Can't remember my password",
      "My login does not work"]

      answer = "If you forgot your My Account password, securely reset it with an authorization code that is sent to your cellphone number on your My Account profile."
    }
  ]
}

resource "awslex_bot_resource" "socal_gas_qnabot" {

  depends_on = [module.bot_sources]
  name       = local.bot_name

  description = local.bot_description

  # path to the bot sources zip file, in bot import/export format
  archive_path = module.bot_sources.archive_path

  # detect changes to the bot sources and update the bot
  source_code_hash = module.bot_sources.archive_sha

  # version of the bot
  alias = "latest"

  # arn of the lambda that fulfills the bot intents
  lambda_arn = local.lambda_arn

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

// give the all aliases associated with this bot permission to invoke the lambda
resource "aws_lambda_permission" "bot_permission" {

  depends_on = [awslex_bot_resource.socal_gas_qnabot]

  statement_id  = "AllowExecutionFromTerraBot"
  action        = "lambda:InvokeFunction"
  function_name = local.lambda_arn
  principal     = "lexv2.amazonaws.com"
  source_arn    = "arn:aws:lex:us-west-2:${local.account_id}:bot-alias/${local.bot_id}/*"
}

output "test_suggestion" {
  depends_on = [awslex_bot_resource.socal_gas_qnabot]
  value      = "aws lexv2-runtime recognize-text --bot-id '${local.bot_id}' --bot-alias-id '${awslex_bot_resource.socal_gas_qnabot.alias_id}' --locale-id 'en_US' --session-id 'test_session' --text 'forgot password'"
}