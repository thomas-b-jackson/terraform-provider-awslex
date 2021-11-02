terraform {
  required_providers {
    awslex = {
      version = "0.1"
      source  = "scg.com/va/awslex"
    }
  }
}

provider "awslex" {
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
  bot_default_alias_id = data.awslex_bot_resource.socal_gas_qnabot.id
}

# create the archive that will be used to create the Lex bot
module "bot_archive" {
  source = "./manifest"
  bot_description = local.bot_description
}

resource "awslex_bot_resource" "socal_gas_qnabot" {

  depends_on = [module.bot_archive]
  name = "TerraBot"

  description = local.bot_description

  # path to the bot manifest archive file, in bot import/export format
  archive_path = module.bot_archive.archive_path

  # version of the bot
  # note: version variable is set to Build.SourceBranch for feature 
  #   branch pipelines, and set to a specific release number on staging or 
  #   prod pipelines
  # note: this results in one, testable bot per feature branch
  alias = "version1"

  # arn of the lambda that fulfills the bot intents
  lambda_arn = local.lambda_versioned_arn

  iam_role = "arn:aws:iam::580753938011:role/aws-service-role/lexv2.amazonaws.com/AWSServiceRoleForLexV2Bots_1PKK306M5NW"
}

// give the bot alias permission to invoke the lambda
resource "aws_lambda_permission" "this_alias" {

  depends_on = [resource.awslex_bot_resource.socal_gas_qnabot]

  statement_id  = "AllowExecutionFromBotAlias"
  action        = "lambda:InvokeFunction"
  function_name = local.lambda_arn
  principal     = "lexv2.amazonaws.com"
  source_arn    = "arn:aws:lex:us-west-2:${local.account_id}:bot-alias/${local.bot_id}/${local.bot_alias_id}"
  qualifier     = local.lambda_version
}

// alias for the bot that always points to the latest version
data "awslex_bot_resource" "socal_gas_qnabot" {
  id = local.bot_id
  alias = "TestBotAlias"
}

// give the bot alias permission to invoke the lambda
resource "aws_lambda_permission" "default_alias" {

  depends_on = [resource.awslex_bot_resource.socal_gas_qnabot,
                data.awslex_bot_resource.socal_gas_qnabot]

  statement_id  = "AllowExecutionFromBotDefaultAlias"
  action        = "lambda:InvokeFunction"
  function_name = local.lambda_arn
  principal     = "lexv2.amazonaws.com"
  source_arn    = "arn:aws:lex:us-west-2:${local.account_id}:bot-alias/${local.bot_id}/${local.bot_default_alias_id}"
  qualifier     = local.lambda_version
}