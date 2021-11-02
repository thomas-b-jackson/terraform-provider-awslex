package provider

import (
	"context"
	"fmt"

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
				Type:        schema.TypeString,
				Required:    true,
				Description: "alias name and version of the bot",
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

	awsClient := meta.(*aws_client.AwsClient)

	err := awsClient.UpdateBot(&bot)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update configured bot",
			Detail:   fmt.Sprintf("Unable to updated configured bot, err: %s", err),
		})
		return diags
	}

	d.SetId(bot.Id)
	// version gets updated with each re-import
	d.Set("version", bot.Version)

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
