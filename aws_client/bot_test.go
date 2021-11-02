package aws_client

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/lexmodelsv2"
	"github.com/aws/aws-sdk-go-v2/service/lexmodelsv2/types"
)

func TestGetBot(t *testing.T) {

	aliasName := "version7"
	botName := "bot-test"
	aliasId := "some-id"
	awsClient, _ := NewTestClient(MockBotClient{
		// outputs from the lexmodelsv2.DescribeBot API call
		DescribeBotOutput: lexmodelsv2.DescribeBotOutput{
			BotName: &botName,
		},
		// outputs from the lexmodelsv2.ListBotAlias API call
		ListBotAliasesOutput: lexmodelsv2.ListBotAliasesOutput{
			BotAliasSummaries: []types.BotAliasSummary{
				{
					BotAliasId:   getAddr("some-other-id"),
					BotAliasName: getAddr("some-other-alias"),
				},
				{
					BotAliasId:   &aliasId,
					BotAliasName: &aliasName,
				},
				{
					BotAliasId:   getAddr("some-some-other-id"),
					BotAliasName: getAddr("some-some-other-alias"),
				},
			},
			BotId: getAddr("testing"),
		},
		// no error (i.e. happy path
		err: nil,
	})

	bot, err := awsClient.GetBot(aliasId, aliasName)

	if err != nil {
		t.Log("error should be nil", err)
		t.Fail()
	}

	fmt.Printf("%+v\n", bot)
}
