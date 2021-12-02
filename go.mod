module github.com/scg/terraform-provider-awslex

go 1.15

require (
	github.com/aws/aws-sdk-go v1.37.0 // indirect
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320
	github.com/hashicorp/hcl/v2 v2.8.2 // indirect
	github.com/hashicorp/terraform-exec v0.15.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.8.0
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/scg/va/aws_client v0.0.0-00010101000000-000000000000
	github.com/zclconf/go-cty v1.10.0 // indirect
	golang.org/x/tools v0.0.0-20201028111035-eafbe7b904eb // indirect
	google.golang.org/api v0.34.0 // indirect
)

replace github.com/scg/va/aws_client => ./aws_client
