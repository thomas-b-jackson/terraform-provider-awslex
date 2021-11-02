module github.com/scg/terraform-provider-awslex

go 1.15

require (
	github.com/aws/aws-sdk-go v1.37.0 // indirect
	github.com/hashicorp/hcl/v2 v2.8.2 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.8.0
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/scg/va/aws_client v0.0.0
	github.com/zclconf/go-cty v1.9.1 // indirect
	golang.org/x/tools v0.0.0-20201028111035-eafbe7b904eb // indirect
	google.golang.org/api v0.34.0 // indirect
)

replace github.com/scg/va/aws_client => ./aws_client
