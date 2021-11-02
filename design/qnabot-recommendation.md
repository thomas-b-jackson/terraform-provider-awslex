## Overview

Options for building v2 Lex Bots resources in Terraform.

## Assumptions

1. We will build a lex QnA bot per [this AWS ML blog](https://aws.amazon.com/blogs/machine-learning/creating-a-question-and-answer-bot-with-amazon-lex-and-amazon-alexa/)
    * per Sempra IAC standards, all resources should ideally be built using terraform via the AWS provider instead of cloud formation
2. Since LexV1 is deprecated, we should use LexV2 
3. Question and answer pairs should be maintained in git
4. Question must be passed as inputs to our terraform skills resource 

## Analysis 

Although the aws terraform provider supports LexV1, it does not yet support LexV2. Also, cloud formation does not yet natively support LexV2. Cloud formation support is likely to arrive soon, but terraform support could lag by up to 1.5 years (based on the 2 year lag between the release of lexv1 and support for lexv2 in the aws provider. lexv2 was released in Jan2021).

With v2 bots, intents and slots are associated with the bot language pack. In other words, intents and slots are attributes of a bot rather than stand-alone resources.

Amazon supports an import/export [schema in json](https://docs.aws.amazon.com/lex/latest/dg/import-export-format.html) for v2 bots, intents and slots. The intended use case for the schema is moving a v2 bot from one amazon account to another, but it could in theory be used to keep a v2 bot in configuration management.

A bot is only build-able/train-able once a version has been assigned to it, and is only testable once an alias has been assigned to the version. 

## Recommendation

To work around the current situation, we recommend one of the following options (in preferred order):
1. create a custom terraform resource that implements a single v2 bot resource (`awslex_bot_resource`)
   * the single resource will serve as an abstraction layer above the sub-resources, including intents, slots, aliases, languages, etc.
   * abandon provider if/when aws provider supports v2 bot resources
2. leverage v2 lex bot cloud formation resources, and call the resource via a `aws_cloudformation_stack` resource in terraform
   * abandon CF if/when aws provider supports v2 bot resources
3. wait for the terraform aws provider to support v2 bot resources
   * based on the v1 bot resources, the provider will likely include separate resources for bots, intents, slots, aliases, etc.
   * wait could be up 18 months from Oct2021

## Preferred Option: awslex_bot_resource resource

Highlights:
* create a `awslex` provider plugin
* implement a `awslex_bot_resource` resource
  * the single resource abstracts sub-resources, resulting in a minimal set of inputs
  * this should dramatically speed plugin development
* the resource uses the aws golang sdk to create the v2 bot, v2 intents, v2 slots, etc.
* the bot, intents, and slots are encoded in json in [v2 bot import/export format](https://docs.aws.amazon.com/lex/latest/dg/import-export-format.html) in the bot repo
* the templated json manifest files are rendered, zipped, and put in an s3 bucket using existing terraform providers
* the manifest file is then imported by the new provider
* the pipeline writes the question/answer pairs into the search index

```
resource "awslex_bot_resource" "socal_gas_qnabot" {

  # path to the archive in s3 containing the bot manifest archive file, 
  # in import/export format
  manifest_s3_path = resource.aws_s3_bucket_object.bot_archive

  # version of the bot
  # note: version variable is set to Build.SourceBranch for feature 
  #   branch pipelines, and set to a specific release number on staging or 
  #   prod pipelines
  # note: this results in one, testable bot per feature branch
  version = var.version

  # arn of the lambda that fulfills the bot intents
  lambda_arn = resource.lambda.fulfillment_lambda.arn
}
```

Bot Build Steps in Provider:
1. import zip from s3 
   * use func (c *Client) StartImport())
   * use existing name
   * assign existing IAM role
   * overwrite existing bot
2. re-build bot 
   * use (func (c *Client) BuildBotLocale())
3. create new version 
   * use (func (c *Client) CreateBotVersion())
4. create new alias 
   * use (func (c *Client) CreateBotAlias())
   * alias name same as version name
   * associate alias with new version
   * associate english language in alias with existing english lambda
5. re-test bot
   
### Pros

* pure terraform
  * full access to terraform Create/Update/Delete operations
* golang-based SDLC
  * fine-grained error handling
  * unit testing

### Cons

* the import/export format is not well documented and could change w/out warning
  * mitigation is straightforward: export the last working version of the bot and update manifest files with export 
* terraform plugins must be written in golang
* bot state not fully reflected in terraform state
  * terraform plan would NOT show changes to individual json making up the bot manifest (would only show changes to the archive zip)

Note: the second con (lo-res tf plans) will be mitigated if/when the aws provider supports lexv2 and the terraform in the bot repo is refactored to leverage the provider.

### Effort

* Dev Effort: 5-7 days
* Pipeline Integration: 1-2 days

Pipeline integration effort includes:
* adding plugin binary

## Second-Best Option: cloud formation + terraform

Highlights:
* leverage AWS::Lex cloud formation resource
* the bot, intents, and slots are encoded in AWS::Lex resources
* the cloud formation template is built using terraform via a `aws_cloudformation_stack` resource

### Pros

* less work than plugin
  
### Cons

* AWS::Lex resources not yet available
  * likely to become available in Nov 2021 (based on conversation with Sempra Amazon TAM Yogesh Chaturvedi in a meeting on Oct26)
  * is likely to be buggy in early releases
* Not compliant with IAC standards
* bot state not reflected in terraform state or reflected in tf plans
  * terraform plan would NOT show changes to individual AWS::Lex resources (would only show a generic change to the `aws_cloudformation_stack` resource)

### Effort

* Dev Effort: 2-3 days
* Pipeline Integration: 2-3 days

Pipeline integration effort includes:
* adding aws secrets
* negotiating an exception for the use of cloud formation

### Third-Best Option: wait for terraform terraform-provider-aws to catch up

Highlights:
* build the bot manually until the terraform aws provider supports lex v2

### Pros

* pure terraform
* bot state fully reflected in terraform state
  
### Cons

* no sense for when the resources will be added to the provider
  * could be 18 months from Oct2021 based on pace of development of lex v1 resources documented on [this PR](https://github.com/hashicorp/terraform-provider-aws/pull/2616)
  * [lex v2 feature request](https://github.com/hashicorp/terraform-provider-aws/issues/21375) is new and has few up-votes or comments

### Effort

* Dev Effort: 2-3 days
* Pipeline Integration: 1 day

Pipeline integration effort includes:
* general debugging