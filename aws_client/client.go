package aws_client

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/lexmodelsv2"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type BotClient interface {
	CreateBot(ctx context.Context, params *lexmodelsv2.CreateBotInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.CreateBotOutput, error)
	UpdateBot(ctx context.Context, params *lexmodelsv2.UpdateBotInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.UpdateBotOutput, error)
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
	ListTagsForResource(ctx context.Context, params *lexmodelsv2.ListTagsForResourceInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.ListTagsForResourceOutput, error)
	TagResource(ctx context.Context, params *lexmodelsv2.TagResourceInput, optFns ...func(*lexmodelsv2.Options)) (*lexmodelsv2.TagResourceOutput, error)
}

// account id and region are needed to form the bot arn
type AwsClient struct {
	Client    BotClient
	AccountId string
	Region    string
}

func NewClient(region string, roleArn string) (*AwsClient, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region))

	if err != nil {
		return nil, err
	}

	var awsClient AwsClient

	if roleArn == "" {

		// assume credentials available in the ~/.aws/credentials file

		// determine account id from sts
		stsClient := sts.NewFromConfig(cfg)
		callerIdentityOutput, accountErr := stsClient.GetCallerIdentity(context.TODO(),
			&sts.GetCallerIdentityInput{})

		if accountErr != nil {
			return nil, accountErr
		}

		client := lexmodelsv2.NewFromConfig(cfg)
		awsClient = AwsClient{client, *callerIdentityOutput.Account, region}

	} else {

		log.Printf("[DEBUG] auth using role arn: %s\n", roleArn)

		// determine account id from sts
		stsClient := sts.NewFromConfig(cfg)
		callerIdentityOutput, accountErr := stsClient.GetCallerIdentity(context.TODO(),
			&sts.GetCallerIdentityInput{})

		if accountErr != nil {
			return nil, accountErr
		}

		// create temporary credentials from the iam role
		stsSvc := sts.NewFromConfig(cfg)
		creds := stscreds.NewAssumeRoleProvider(stsSvc, roleArn)
		cfg.Credentials = aws.NewCredentialsCache(creds)
		client := lexmodelsv2.NewFromConfig(cfg)

		awsClient = AwsClient{client, *callerIdentityOutput.Account, region}
	}

	return &awsClient, nil
}
