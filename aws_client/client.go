package aws_client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lexmodelsv2"
)

type BotClient interface {
	CreateBot(ctx context.Context, params *lexmodelsv2.CreateBotInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.CreateBotOutput, error)
	DescribeBot(ctx context.Context, params *lexmodelsv2.DescribeBotInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.DescribeBotOutput, error)
	ListBotAliases(ctx context.Context, params *lexmodelsv2.ListBotAliasesInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.ListBotAliasesOutput, error)
	CreateBotAlias(ctx context.Context, params *lexmodelsv2.CreateBotAliasInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.CreateBotAliasOutput, error)
	CreateBotVersion(ctx context.Context, params *lexmodelsv2.CreateBotVersionInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.CreateBotVersionOutput, error)
	DescribeBotVersion(ctx context.Context, params *lexmodelsv2.DescribeBotVersionInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.DescribeBotVersionOutput, error)
	DescribeBotAlias(ctx context.Context, params *lexmodelsv2.DescribeBotAliasInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.DescribeBotAliasOutput, error)
	CreateUploadUrl(ctx context.Context, params *lexmodelsv2.CreateUploadUrlInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.CreateUploadUrlOutput, error)
	StartImport(ctx context.Context, params *lexmodelsv2.StartImportInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.StartImportOutput, error)
	DescribeImport(ctx context.Context, params *lexmodelsv2.DescribeImportInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.DescribeImportOutput, error)
	DeleteBot(ctx context.Context, params *lexmodelsv2.DeleteBotInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.DeleteBotOutput, error)
	ListBotVersions(ctx context.Context, params *lexmodelsv2.ListBotVersionsInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.ListBotVersionsOutput, error)
	UpdateBotAlias(ctx context.Context, params *lexmodelsv2.UpdateBotAliasInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.UpdateBotAliasOutput, error)
	BuildBotLocale(ctx context.Context, params *lexmodelsv2.BuildBotLocaleInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.BuildBotLocaleOutput, error)
	DescribeBotLocale(ctx context.Context, params *lexmodelsv2.DescribeBotLocaleInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.DescribeBotLocaleOutput, error)
}

type AwsClient struct {
	Client BotClient
}

func NewClient() (*AwsClient, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := lexmodelsv2.NewFromConfig(cfg)

	c := AwsClient{client}

	return &c, nil
}
