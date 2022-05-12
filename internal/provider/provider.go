package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scg/va/aws_client"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown
}

func Provider(version string) *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AWS_DEFAULT_REGION", nil),
			},
			"role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"awslex_bot_resource": dataSourceBot(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"awslex_bot_resource": resourceBot(),
		},
		ConfigureContextFunc: configure,
	}
}

func configure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	var diags diag.Diagnostics

	region := d.Get("region").(string)
	role_arn := d.Get("role_arn").(string)

	askClient, err := aws_client.NewClient(region, role_arn)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create aws client",
			Detail:   fmt.Sprintf("Unable to create aws client: %s", err),
		})
		return nil, diags
	}

	return askClient, diags
}
