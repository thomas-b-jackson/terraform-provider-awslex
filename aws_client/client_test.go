package aws_client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/lexmodelsv2"
)

// allow tests to pass in their own mock client
func NewTestClient(client BotClient) (*AwsClient, error) {
	c := AwsClient{client}
	return &c, nil
}

// for use in unit tests, were we want to allow the test to specify
// the outputs of a given api call
type MockBotClient struct {
	BotClient
	// each test should specify the expected output and error
	DescribeBotOutput      lexmodelsv2.DescribeBotOutput
	ListBotAliasesOutput   lexmodelsv2.ListBotAliasesOutput
	DescribeBotAliasOutput lexmodelsv2.DescribeBotAliasOutput
	err                    error
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
