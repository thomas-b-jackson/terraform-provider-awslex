package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scg/va/aws_client"
)

func resourceBot() *schema.Resource {
	return &schema.Resource{
		Description: "Alex skill resource",

		CreateContext: resourceBotCreate,
		ReadContext:   resourceBotRead,
		UpdateContext: resourceBotUpdate,
		DeleteContext: resourceBotDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the bot",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the bot",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of the bot",
			},
			"alias": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "alias name and version of the bot",
				ValidateDiagFunc: AliasValidator,
			},
			"alias_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the bot alias",
			},
			"archive_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Path to the zip archive containing intents and slots",
			},
			"source_code_hash": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Base64-encoded representation of the SHA-256 sum of the zip file",
			},
			"lambda_arn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Arn of router lambda",
			},
			"iam_role": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Arn of IAM role to use with the bot",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Description of bot",
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceBotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	var bot aws_client.LexBot

	bot.Name = d.Get("name").(string)
	bot.Alias = d.Get("alias").(string)
	bot.LambdaArn = d.Get("lambda_arn").(string)
	bot.Description = d.Get("description").(string)
	bot.IamRoleArn = d.Get("iam_role").(string)
	bot.ArchivePath = d.Get("archive_path").(string)
	bot.SourceCodeHash = d.Get("source_code_hash").(string)
	bot.Tags = convertTags(d.Get("tags").(map[string]interface{}))

	awsClient := meta.(*aws_client.AwsClient)

	err := awsClient.CreateBot(&bot)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create configured bot",
			Detail:   fmt.Sprintf("Unable to create configured bot, err: %s", err),
		})
		return diags
	}

	// set computed values
	d.SetId(bot.Id)
	d.Set("version", bot.Version)
	d.Set("alias_id", bot.AliasId)

	return diags
}

func convertTags(tags map[string]interface{}) map[string]string {
	result := map[string]string{}
	for k := range tags {
		if v, ok := tags[k].(string); ok {
			result[k] = v
		}
	}
	return result
}

func resourceBotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return dataSourceBotRead(ctx, d, meta)
}

func resourceBotUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	var bot aws_client.LexBot

	bot.Id = d.Get("id").(string)
	bot.Name = d.Get("name").(string)
	bot.Alias = d.Get("alias").(string)
	bot.AliasId = d.Get("alias_id").(string)
	bot.LambdaArn = d.Get("lambda_arn").(string)
	bot.Description = d.Get("description").(string)
	bot.IamRoleArn = d.Get("iam_role").(string)
	bot.ArchivePath = d.Get("archive_path").(string)
	bot.Version = d.Get("version").(string)
	bot.SourceCodeHash = d.Get("source_code_hash").(string)
	bot.Tags = convertTags(d.Get("tags").(map[string]interface{}))

	awsClient := meta.(*aws_client.AwsClient)

	err := awsClient.UpdateBot(&bot, d)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update configured bot",
			Detail:   fmt.Sprintf("Unable to updated configured bot, err: %s", err),
		})
		return diags
	}

	d.SetId(bot.Id)
	// version gets updated with each update
	d.Set("version", bot.Version)
	// alias id may get updated with each update
	d.Set("alias_id", bot.AliasId)

	return diags
}

func resourceBotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	var bot aws_client.LexBot

	// TODO: remove hardcoding later
	botId := d.Get("id").(string)

	awsClient := meta.(*aws_client.AwsClient)

	err := awsClient.DeleteBot(botId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create requested bot",
			Detail:   fmt.Sprintf("Unable to create requested bot, err: %s", err),
		})
		return diags
	}

	d.SetId(bot.Id)

	return diags
}

func AliasValidator(i interface{}, p cty.Path) diag.Diagnostics {
	alias := i.(string)

	match, err := regexp.Match("^[A-Za-z0-9_-]+$", []byte(alias))

	// log.Printf("[DEBUG] regex match? match: %t, err: %s\n", match, err)

	if err != nil || !match {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Alias contains invalid characters",
				Detail:   "Alias contains invalid characters. Valid characters: A-Z, a-z, 0-9, -, _",
			},
		}
	}

	return diag.Diagnostics{}
}
