terraform {
  required_providers {
    awslex = {
      version = "0.1"
      source  = "scg.com/va/awslex"
    }
  }
}

provider "awslex" {
  // region = "us-west-2"
}

provider "aws" {
  region = "us-west-2"
}

data "aws_caller_identity" "current" {}

locals {
  bot_description = "Terraform Bot"
  lambda_version = 1
  lambda_arn = "arn:aws:lambda:us-west-2:580753938011:function:QnABot-FulfillmentLambda-iGXdhe8RHdyH"
  lambda_versioned_arn = "${local.lambda_arn}:${local.lambda_version}"
  account_id = data.aws_caller_identity.current.account_id
  bot_id = resource.awslex_bot_resource.socal_gas_qnabot.id
  bot_alias_id = resource.awslex_bot_resource.socal_gas_qnabot.alias_id
}  

# create the file that represents the Lex bot sources
module "bot_sources" {
  source = "./sources"
  bot_description = local.bot_description
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
    },
    {
        id = "pilot-light"
        questions = ["what should i do if my pilot light is out?",
                    "pilot light",
                    "help with pilot light on furnace",
                    "pilot light on furnace"]

        answer = "If you have a gas water heater or gas furnace that is not working, the pilot light on the heater may have accidentally become extinguished. For assistance relighting the pilot light, consult the water heater or furnace owner’s manual or contact Southern California Gas Company (SoCalGas®) and schedule an appliance appointment."
    }    
  ]
}

resource "awslex_bot_resource" "socal_gas_qnabot" {

  depends_on = [module.bot_sources]
  name = "TerraBot"

  description = local.bot_description

  # path to the bot sources zip file, in bot import/export format
  archive_path = module.bot_sources.archive_path

  # detect changes to the bot sources and update the bot
  source_code_hash = module.bot_sources.archive_sha

  # version of the bot
  # note: version variable should be set to Build.SourceBranch for feature 
  #   branch pipelines, and set to a specific release number on staging or 
  #   prod pipelines
  # note: this results in one, testable bot per feature branch
  alias = "fix_gas-leaks"

  # arn of the lambda that fulfills the bot intents
  lambda_arn = local.lambda_versioned_arn

  iam_role = "arn:aws:iam::580753938011:role/aws-service-role/lexv2.amazonaws.com/AWSServiceRoleForLexV2Bots_1PKK306M5NW"
}

// give the all aliases associated with this bot permission to invoke the lambda
resource "aws_lambda_permission" "bot_permission" {

  depends_on = [resource.awslex_bot_resource.socal_gas_qnabot]

  statement_id  = "AllowExecutionFromBotAlias"
  action        = "lambda:InvokeFunction"
  function_name = local.lambda_arn
  principal     = "lexv2.amazonaws.com"
  source_arn    = "arn:aws:lex:us-west-2:${local.account_id}:bot-alias/${local.bot_id}/*"
  qualifier     = local.lambda_version
}

output "test_suggestion" {
  depends_on = [resource.awslex_bot_resource.socal_gas_qnabot]
  value = "aws lexv2-runtime recognize-text --bot-id '${local.bot_id}' --bot-alias-id '${resource.awslex_bot_resource.socal_gas_qnabot.alias_id}' --locale-id 'en_US' --session-id 'test_session' --text 'forgot password'"
}