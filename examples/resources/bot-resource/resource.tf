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

resource "awslex_bot_resource" "socal_gas_qnabot" {

  name = "TerraBot"

  description = "Terraform Bot"
  # path to the archive in s3 containing the bot manifest archive file, 
  # in import/export format
  archive_path = "/mnt/c/Users/tomj/Downloads/QnABot_QnaBot-6-USX2IJSEYW-LexJson.zip"
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

data "aws_caller_identity" "current" {}

locals {
  lambda_version = 1
  lambda_arn = "arn:aws:lambda:us-west-2:580753938011:function:QnABot-FulfillmentLambda-iGXdhe8RHdyH"
  lambda_versioned_arn = "${local.lambda_arn}:${local.lambda_version}"
  account_id = data.aws_caller_identity.current.account_id
  bot_id = resource.awslex_bot_resource.socal_gas_qnabot.id
  bot_alias_id = resource.awslex_bot_resource.socal_gas_qnabot.alias_id
  bot_default_alias_id = data.awslex_bot_resource.socal_gas_qnabot.id
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