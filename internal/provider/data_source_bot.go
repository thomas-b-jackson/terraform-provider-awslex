package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scg/va/aws_client"
)

func dataSourceBot() *schema.Resource {

	return &schema.Resource{
		Description: "a data source returning details on a v2 lex bot",
		ReadContext: dataSourceBotRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the bot",
			},
			"alias": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Alias name of the bot",
				ValidateDiagFunc: AliasValidator,
			},
			"alias_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the bot alias",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the bot",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of bot",
			},
			"lambda_arn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Arn of router lambda",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of bot",
			},
			"iam_role": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAM role of bot",
			},
			"source_code_hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Base64-encoded representation of raw SHA-256 sum of the zip file",
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceBotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	botId := d.Get("id").(string)
	botAlias := d.Get("alias").(string)

	awsClient := meta.(*aws_client.AwsClient)

	bot, err := awsClient.GetBot(botId, botAlias)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get requested bot",
			Detail:   fmt.Sprintf("Unable to get requested bot, err: %s", err),
		})
		return diags
	}

	// set computed values
	d.SetId(botId)
	d.Set("lambda_arn", bot.LambdaArn)
	d.Set("name", bot.Name)
	d.Set("iam_role", bot.IamRoleArn)
	d.Set("description", bot.Description)
	d.Set("version", bot.Version)
	d.Set("alias_id", bot.AliasId)
	d.Set("source_code_hash", bot.SourceCodeHash)
	d.Set("tags", bot.Tags)
	return diags
}
