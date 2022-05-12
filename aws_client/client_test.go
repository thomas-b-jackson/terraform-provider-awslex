package aws_client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/lexmodelsv2"
)

// allow tests to pass in their own mock client
func NewTestClient(client BotClient) (*AwsClient, error) {
	c := AwsClient{client, "abcd", "us-west-2"}
	return &c, nil
}

// for use in unit tests, were we want to allow the test to specify
// the outputs of a given api call
type MockBotClient struct {
	BotClient
	// each test should specify the expected output and error
	DescribeBotOutput         lexmodelsv2.DescribeBotOutput
	ListBotAliasesOutput      lexmodelsv2.ListBotAliasesOutput
	DescribeBotAliasOutput    lexmodelsv2.DescribeBotAliasOutput
	DescribeBotVersionOutput  lexmodelsv2.DescribeBotVersionOutput
	ListTagsForResourceOutput lexmodelsv2.ListTagsForResourceOutput
	TagResourceOutput         lexmodelsv2.TagResourceOutput
	err                       error
}

func (m MockBotClient) ListBotAliases(ctx context.Context, params *lexmodelsv2.ListBotAliasesInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.ListBotAliasesOutput, error) {
	return &m.ListBotAliasesOutput, m.err
}
func (m MockBotClient) DescribeBotAlias(ctx context.Context, params *lexmodelsv2.DescribeBotAliasInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.DescribeBotAliasOutput, error) {
	return &m.DescribeBotAliasOutput, m.err
}
func (m MockBotClient) DescribeBot(ctx context.Context, params *lexmodelsv2.DescribeBotInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.DescribeBotOutput, error) {
	return &m.DescribeBotOutput, m.err
}
func (m MockBotClient) DescribeBotVersion(ctx context.Context, params *lexmodelsv2.DescribeBotVersionInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.DescribeBotVersionOutput, error) {
	return &m.DescribeBotVersionOutput, m.err
}
func (m MockBotClient) ListTagsForResource(ctx context.Context, params *lexmodelsv2.ListTagsForResourceInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.ListTagsForResourceOutput, error) {
	return &m.ListTagsForResourceOutput, m.err
}
func (m MockBotClient) TagResource(ctx context.Context, params *lexmodelsv2.TagResourceInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.TagResourceOutput, error) {
	return &m.TagResourceOutput, m.err
}
