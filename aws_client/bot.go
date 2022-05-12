package aws_client

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/lexmodelsv2"
	"github.com/aws/aws-sdk-go-v2/service/lexmodelsv2/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type LexBot struct {
	Id             string
	Name           string
	Alias          string
	AliasId        string
	Version        string
	ArchivePath    string
	Description    string
	LambdaArn      string
	IamRoleArn     string
	SourceCodeHash string
	Tags           map[string]string
}

// wait up to this many seconds for long-running bot operations to to complete
const BotWaitTimeoutSec = 60
const DraftVersion = "DRAFT"

var ttl int32 = 100

func (c *AwsClient) GetBot(botId string, alias string) (LexBot, error) {

	var bot LexBot

	botDescription, err := c.Client.DescribeBot(context.TODO(),
		&lexmodelsv2.DescribeBotInput{
			BotId: &botId,
		})

	if err != nil {
		return LexBot{}, err
	}

	bot.Name = *botDescription.BotName
	if botDescription.Description != nil {
		bot.Description = *botDescription.Description
	}
	bot.IamRoleArn = *botDescription.RoleArn

	botAlias, err := c.Client.ListBotAliases(context.TODO(),
		&lexmodelsv2.ListBotAliasesInput{
			BotId: &botId,
		})

	if err != nil {
		return bot, err
	}

	for _, botAlias := range botAlias.BotAliasSummaries {
		if *botAlias.BotAliasName == alias {

			bot.Id = botId
			bot.Alias = *botAlias.BotAliasName
			bot.AliasId = *botAlias.BotAliasId
			// not all aliases have a version
			if botAlias.BotVersion != nil {
				bot.Version = *botAlias.BotVersion
			}
		}
	}

	if bot.AliasId != "" {

		// get tags associated with the bot alias
		listTagsForResourceOutput, err := c.Client.ListTagsForResource(context.TODO(),
			&lexmodelsv2.ListTagsForResourceInput{
				ResourceARN: getAddr(getAliasArn(botId, bot.AliasId, c.Region, c.AccountId)),
			})

		if err == nil && listTagsForResourceOutput != nil && listTagsForResourceOutput.Tags != nil {
			bot.Tags = make(map[string]string)
			for key, val := range listTagsForResourceOutput.Tags {
				bot.Tags[key] = val
			}
		}

		// describe the bot to get its lambda arn
		var describeBotAliasOutput *lexmodelsv2.DescribeBotAliasOutput

		describeBotAliasOutput, err = c.Client.DescribeBotAlias(context.TODO(),
			&lexmodelsv2.DescribeBotAliasInput{
				BotAliasId: &bot.AliasId,
				BotId:      &botId,
			})

		if err != nil {
			return LexBot{}, fmt.Errorf("error describing bot alias %s: %s", bot.AliasId, err)
		}

		if describeBotAliasOutput.BotAliasLocaleSettings != nil {
			usLocalSettings, ok := describeBotAliasOutput.BotAliasLocaleSettings["en_US"]

			if ok && usLocalSettings.CodeHookSpecification != nil {
				bot.LambdaArn = *usLocalSettings.CodeHookSpecification.LambdaCodeHook.LambdaARN
			}
		}
	}

	if bot.Version != "" && bot.Version != DraftVersion {

		// describe the bot version to get its source code hash
		var describeBotVersionOutput *lexmodelsv2.DescribeBotVersionOutput

		describeBotVersionOutput, err = c.Client.DescribeBotVersion(context.TODO(),
			&lexmodelsv2.DescribeBotVersionInput{
				BotId:      &botId,
				BotVersion: &bot.Version,
			})

		if err != nil {
			return LexBot{}, fmt.Errorf("error describing bot version %s. cannot get hash: %s", bot.AliasId, err)
		}

		// the description field is used to store the source code hash
		if describeBotVersionOutput.Description != nil {
			bot.SourceCodeHash = *describeBotVersionOutput.Description
		}
	}

	return bot, err
}

func (c *AwsClient) CreateBot(bot *LexBot) error {

	// create the bot skeleton in aws
	err := c.createBot(bot)
	if err != nil {
		return err
	}

	// put the archive containing intents and slots in s3
	// (in a location determined by the aws lex sdk)
	uploadId, err := c.upload(bot.ArchivePath)

	if err != nil {
		return err
	}

	// import the bot intents and slots into the bot
	err = c.importBot(uploadId, *bot)

	if err != nil {
		return err
	}

	// set the version of the imported bot
	err = c.setImportedVersion(bot)

	if err != nil {
		return err
	}

	// update the original alias to reference the desired lambda
	err = c.updateOriginalAlias(bot)

	if err != nil {
		return err
	}

	// build the bot
	err = c.buildBot(bot)

	if err != nil {
		return err
	}

	// create a new version for the imported bot
	err = c.createVersion(bot)

	if err != nil {
		return err
	}

	// create an alias to the new version whose name matches the
	// alias defined in the tf bot resource
	err = c.createAlias(bot)

	if err != nil {
		return err
	}

	return err
}

func (c *AwsClient) UpdateBot(bot *LexBot, d *schema.ResourceData) error {

	var err error

	if d.HasChange("name") || d.HasChange("description") || d.HasChange("iam_role") {

		_, err := c.Client.UpdateBot(context.TODO(), &lexmodelsv2.UpdateBotInput{
			BotId:   &bot.Id,
			BotName: &bot.Name,
			DataPrivacy: &types.DataPrivacy{
				ChildDirected: false,
			},
			RoleArn:                 &bot.IamRoleArn,
			Description:             &bot.Description,
			IdleSessionTTLInSeconds: &ttl,
		})

		if err != nil {
			return err
		}
	}

	if d.HasChange("source_code_hash") {

		// put the archive containing intents and slots in s3
		// (in a location determined by the aws lex sdk)
		uploadId, err := c.upload(bot.ArchivePath)

		if err != nil {
			return err
		}

		// import the bot intents and slots into the bot
		err = c.importBot(uploadId, *bot)

		if err != nil {
			return err
		}

		// create a new version for the imported bot
		err = c.createVersion(bot)

		if err != nil {
			return err
		}

		// build the bot
		err = c.buildBot(bot)
	}

	// create or update alias for the bot
	err = c.createOrUpdateAlias(bot)

	if err != nil {
		return err
	}

	if d.HasChange("tags") {
		// updated tags on alias
		log.Printf("[DEBUG] adding tags to alias: %v\n", bot.Tags)
		_, err = c.Client.TagResource(context.TODO(), &lexmodelsv2.TagResourceInput{
			ResourceARN: getAddr(getAliasArn(bot.Id, bot.AliasId, c.Region, c.AccountId)),
			Tags:        bot.Tags,
		})
	}

	return err
}
func (c *AwsClient) createBot(bot *LexBot) error {

	createBotOutput, err := c.Client.CreateBot(context.TODO(), &lexmodelsv2.CreateBotInput{
		BotName: &bot.Name,
		DataPrivacy: &types.DataPrivacy{
			ChildDirected: false,
		},
		RoleArn:                 &bot.IamRoleArn,
		Description:             &bot.Description,
		IdleSessionTTLInSeconds: &ttl,
	})

	if err != nil {
		return err
	}

	bot.Id = *createBotOutput.BotId

	// wait for creation to complete
	expiredTimeSec := 0
	sleepDurationSec := 10
	for {
		botDescription, describeErr := c.Client.DescribeBot(context.TODO(),
			&lexmodelsv2.DescribeBotInput{
				BotId: &bot.Id,
			})

		// break if creation is complete
		if (describeErr == nil && botDescription.BotStatus == types.BotStatusAvailable) ||
			expiredTimeSec >= BotWaitTimeoutSec {
			break
		}

		log.Printf("[DEBUG] waiting for bot creation to complete. Current status: %s\n", botDescription.BotStatus)

		// sleep for X seconds
		time.Sleep(time.Duration(sleepDurationSec) * time.Second)
		expiredTimeSec += sleepDurationSec
	}

	return err
}

func (c *AwsClient) importBot(uploadId string,
	bot LexBot) error {

	// import the archive
	_, err := c.Client.StartImport(context.TODO(), &lexmodelsv2.StartImportInput{
		ImportId:      &uploadId,
		MergeStrategy: types.MergeStrategyOverwrite,
		ResourceSpecification: &types.ImportResourceSpecification{
			BotImportSpecification: &types.BotImportSpecification{
				BotName: &bot.Name,
				DataPrivacy: &types.DataPrivacy{
					ChildDirected: false,
				},
				RoleArn: &bot.IamRoleArn,
			},
		},
	})

	if err != nil {
		return err
	}

	// wait for import to complete
	expiredTimeSec := 0
	sleepDurationSec := 10
	for {
		startImportOutput, err := c.Client.DescribeImport(context.TODO(), &lexmodelsv2.DescribeImportInput{
			ImportId: &uploadId,
		})

		// log.Printf("[DEBUG] import output: %v\n", startImportOutput)
		// log.Printf("[DEBUG] import status: %v\n", startImportOutput.ImportStatus)

		// break if import is complete
		if (err == nil && startImportOutput.ImportStatus == types.ImportStatusCompleted) ||
			expiredTimeSec >= BotWaitTimeoutSec {
			break
		}

		log.Printf("[DEBUG] waiting for bot import to complete. Current status: %s\n", startImportOutput.ImportStatus)

		// sleep for X seconds
		time.Sleep(time.Duration(sleepDurationSec) * time.Second)
		expiredTimeSec += sleepDurationSec
	}

	return err
}

func (c *AwsClient) upload(archivePath string) (string, error) {

	var uploadId string

	createUploadUrlOutput, err := c.Client.CreateUploadUrl(context.TODO(), &lexmodelsv2.CreateUploadUrlInput{})

	if err != nil {
		return uploadId, err
	}

	uploadId = *createUploadUrlOutput.ImportId
	uploadUrl := *createUploadUrlOutput.UploadUrl

	// log.Printf("[DEBUG] create url id: %s, url: %s\n", *createUploadUrlOutput.ImportId, *createUploadUrlOutput.UploadUrl)

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	b, err := ioutil.ReadFile(archivePath)
	if err != nil {
		return uploadId, err
	}

	req, err := http.NewRequest("PUT", uploadUrl, bytes.NewReader(b))
	if err != nil {
		return uploadId, err
	}

	contentType := http.DetectContentType(b)
	req.Header.Set("Content-Type", contentType)
	rsp, _ := client.Do(req)

	// log.Printf("[DEBUG] upload post content type %v\n", contentType)
	// log.Printf("[DEBUG] upload post response %v\n", rsp)

	if rsp.StatusCode != http.StatusOK {
		return uploadId, fmt.Errorf("upload failed with response code: %d", rsp.StatusCode)
	}
	return uploadId, nil
}

func (c *AwsClient) setImportedVersion(bot *LexBot) error {

	listBotVersionOutput, err := c.Client.ListBotVersions(context.TODO(), &lexmodelsv2.ListBotVersionsInput{
		BotId: &bot.Id,
	})

	if err != nil {
		return err
	}

	if listBotVersionOutput != nil && len(listBotVersionOutput.BotVersionSummaries) == 1 {
		bot.Version = *listBotVersionOutput.BotVersionSummaries[0].BotVersion
	} else {
		err = fmt.Errorf("could not find bot version for imported bot: %s", bot.Id)
	}

	return err
}

func (c *AwsClient) updateOriginalAlias(bot *LexBot) error {

	listBotAliasOutput, err := c.Client.ListBotAliases(context.TODO(), &lexmodelsv2.ListBotAliasesInput{
		BotId: &bot.Id,
	})

	if err != nil {
		return err
	}

	var ogAliasId string
	var ogAliasName string

	if listBotAliasOutput != nil && len(listBotAliasOutput.BotAliasSummaries) == 1 {
		ogAliasId = *listBotAliasOutput.BotAliasSummaries[0].BotAliasId
		ogAliasName = *listBotAliasOutput.BotAliasSummaries[0].BotAliasName
	} else {
		err = fmt.Errorf("could not find bot alias for imported bot: %s", bot.Id)
		return err
	}

	// update the alias to point to the lambda function
	_, err = c.Client.UpdateBotAlias(context.TODO(), &lexmodelsv2.UpdateBotAliasInput{
		BotId:        &bot.Id,
		BotAliasId:   &ogAliasId,
		BotAliasName: &ogAliasName,
		BotVersion:   &bot.Version,
		BotAliasLocaleSettings: map[string]types.BotAliasLocaleSettings{
			"en_US": {
				CodeHookSpecification: &types.CodeHookSpecification{
					LambdaCodeHook: &types.LambdaCodeHook{
						LambdaARN:                &bot.LambdaArn,
						CodeHookInterfaceVersion: getAddr("1.0"),
					},
				},
				Enabled: true,
			},
		},
	})

	return err
}

func (c *AwsClient) createVersion(bot *LexBot) error {

	createBotVersionOutput, err := c.Client.CreateBotVersion(context.TODO(), &lexmodelsv2.CreateBotVersionInput{
		BotId: &bot.Id,
		// use the description field to store the source code hash
		Description: &bot.SourceCodeHash,
		BotVersionLocaleSpecification: map[string]types.BotVersionLocaleDetails{
			"en_US": {
				SourceBotVersion: getAddr(bot.Version),
			},
		},
	})

	if err != nil {
		return err
	}

	bot.Version = *createBotVersionOutput.BotVersion

	// wait for version to become available
	expiredTimeSec := 0
	sleepDurationSec := 10
	for {
		describeBotVersionOutput, err := c.Client.DescribeBotVersion(context.TODO(), &lexmodelsv2.DescribeBotVersionInput{
			BotId:      &bot.Id,
			BotVersion: &bot.Version,
		})

		// break if version is available
		if (err == nil && describeBotVersionOutput.BotStatus == types.BotStatusAvailable) ||
			expiredTimeSec >= BotWaitTimeoutSec {
			break
		}

		if describeBotVersionOutput != nil {
			log.Printf("[DEBUG] waiting for bot version to become available. Current status: %s\n", describeBotVersionOutput.BotStatus)
		} else {
			log.Printf("[DEBUG] waiting for bot version to become available. Current status: %s\n", "unknown")
		}

		// sleep for X seconds
		time.Sleep(time.Duration(sleepDurationSec) * time.Second)
		expiredTimeSec += sleepDurationSec
	}

	return err
}

func (c *AwsClient) buildBot(bot *LexBot) error {

	_, err := c.Client.BuildBotLocale(context.TODO(), &lexmodelsv2.BuildBotLocaleInput{
		BotId: &bot.Id,
		// The version of the bot to build can only be the draft version
		BotVersion: getAddr(DraftVersion),
		LocaleId:   getAddr("en_US"),
	})

	if err != nil {
		return err
	}

	// wait for build to complete
	expiredTimeSec := 0
	sleepDurationSec := 10
	for {
		describeBotLocaleOutput, err := c.Client.DescribeBotLocale(context.TODO(), &lexmodelsv2.DescribeBotLocaleInput{
			BotId:      &bot.Id,
			BotVersion: getAddr(DraftVersion),
			LocaleId:   getAddr("en_US"),
		})

		// break if version is available
		if (err == nil && describeBotLocaleOutput.BotLocaleStatus == types.BotLocaleStatusBuilt) ||
			expiredTimeSec >= BotWaitTimeoutSec {
			break
		}

		if describeBotLocaleOutput != nil {
			log.Printf("[DEBUG] waiting for bot build to complete. Current status: %s\n", describeBotLocaleOutput.BotLocaleStatus)
		} else {
			log.Printf("[DEBUG] waiting for bot build to complete. Current status: %s\n", "unknown")
		}

		// sleep for X seconds
		time.Sleep(time.Duration(sleepDurationSec) * time.Second)
		expiredTimeSec += sleepDurationSec
	}

	return err
}

func (c *AwsClient) createOrUpdateAlias(bot *LexBot) error {

	// see if the alias already exists
	aliasId, err := c.getAliasId(bot, bot.Alias)

	if err != nil {
		return err
	}

	if aliasId == "" {
		err = c.createAlias(bot)
	} else {
		bot.AliasId = aliasId
		err = c.updateAlias(bot)
	}

	return err
}

func (c *AwsClient) updateAlias(bot *LexBot) error {

	// update the existing alias to reference the bot version
	_, err := c.Client.UpdateBotAlias(context.TODO(), &lexmodelsv2.UpdateBotAliasInput{
		BotId:        &bot.Id,
		BotAliasId:   &bot.AliasId,
		BotAliasName: &bot.Alias,
		BotVersion:   &bot.Version,
		BotAliasLocaleSettings: map[string]types.BotAliasLocaleSettings{
			"en_US": {
				CodeHookSpecification: &types.CodeHookSpecification{
					LambdaCodeHook: &types.LambdaCodeHook{
						LambdaARN:                &bot.LambdaArn,
						CodeHookInterfaceVersion: getAddr("1.0"),
					},
				},
				Enabled: true,
			},
		},
	})

	if err != nil {
		return err
	}

	// wait for the alias to become available
	c.aliasWait(bot)

	return err
}

func (c *AwsClient) createAlias(bot *LexBot) error {

	botTags := make(map[string]string)
	for key, val := range bot.Tags {
		botTags[key] = val
	}

	// create the alias
	createBotAliasOutput, err := c.Client.CreateBotAlias(context.TODO(), &lexmodelsv2.CreateBotAliasInput{
		BotId:        &bot.Id,
		BotAliasName: &bot.Alias,
		BotVersion:   &bot.Version,
		Tags:         botTags,
		BotAliasLocaleSettings: map[string]types.BotAliasLocaleSettings{
			"en_US": {
				CodeHookSpecification: &types.CodeHookSpecification{
					LambdaCodeHook: &types.LambdaCodeHook{
						LambdaARN:                &bot.LambdaArn,
						CodeHookInterfaceVersion: getAddr("1.0"),
					},
				},
				Enabled: true,
			},
		},
	})

	if err != nil {
		return err
	}

	// save the id of the bot alias
	bot.AliasId = *createBotAliasOutput.BotAliasId

	// wait for the alias to become available
	c.aliasWait(bot)

	return err
}

func (c *AwsClient) aliasWait(bot *LexBot) {
	// wait for bot alias to be available
	expiredTimeSec := 0
	sleepDurationSec := 10
	for {
		describeBotAliasOutput, describeErr := c.Client.DescribeBotAlias(context.TODO(),
			&lexmodelsv2.DescribeBotAliasInput{
				BotId:      &bot.Id,
				BotAliasId: &bot.AliasId,
			})

		// break if alias is available
		if (describeErr == nil && describeBotAliasOutput.BotAliasStatus == types.BotAliasStatusAvailable) ||
			expiredTimeSec >= BotWaitTimeoutSec {
			break
		}

		if describeBotAliasOutput != nil {
			log.Printf("[DEBUG] waiting for bot alias to be available. Current status: %s\n", describeBotAliasOutput.BotAliasStatus)
		} else {
			log.Printf("[DEBUG] waiting for bot alias to be available. Current status: %s\n", "unknown")
		}

		// sleep for X seconds
		time.Sleep(time.Duration(sleepDurationSec) * time.Second)
		expiredTimeSec += sleepDurationSec
	}
}

func (c *AwsClient) getAliasId(bot *LexBot, alias string) (string, error) {
	botAlias, err := c.Client.ListBotAliases(context.TODO(),
		&lexmodelsv2.ListBotAliasesInput{
			BotId: &bot.Id,
		})

	if err != nil {
		return "", err
	}

	for _, botAlias := range botAlias.BotAliasSummaries {
		if *botAlias.BotAliasName == alias {
			return *botAlias.BotAliasId, err
		}
	}
	return "", err
}

func (c *AwsClient) DeleteBot(botId string) error {

	_, err := c.Client.DeleteBot(context.TODO(), &lexmodelsv2.DeleteBotInput{
		BotId:                  &botId,
		SkipResourceInUseCheck: true,
	})

	if err != nil {
		return err
	}

	// wait for deletion to complete
	expiredTimeSec := 0
	sleepDurationSec := 10
	for {
		botDescription, describeErr := c.Client.DescribeBot(context.TODO(),
			&lexmodelsv2.DescribeBotInput{
				BotId: &botId,
			})

		// break if deletion is complete
		if (describeErr == nil && botDescription.BotStatus != types.BotStatusDeleting) ||
			expiredTimeSec >= BotWaitTimeoutSec {
			break
		}

		if botDescription != nil {
			log.Printf("[DEBUG] waiting for bot deletion to complete. Current status: %s\n", botDescription.BotStatus)
		} else {
			// assume deletion is complete if bot description is not available
			break
		}

		// sleep for X seconds
		time.Sleep(time.Duration(sleepDurationSec) * time.Second)
		expiredTimeSec += sleepDurationSec
	}

	return err
}

func getAddr(s string) *string {
	return &s
}

func getAliasArn(botId string, aliasId string, accountId string, region string) string {
	return fmt.Sprintf("arn:aws:lex:%s:%s:bot-alias/%s/%s", accountId, region, botId, aliasId)
}
