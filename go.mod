module github.com/scg/terraform-provider-awslex

go 1.15

require (
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.16.0
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/scg/va/aws_client v0.0.0-00010101000000-000000000000
	google.golang.org/genproto v0.0.0-20200904004341-0bd0a958aa1d // indirect
)

replace github.com/scg/va/aws_client => ./aws_client
