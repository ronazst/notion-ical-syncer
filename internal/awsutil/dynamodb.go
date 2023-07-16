package awsutil

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/ronazst/notion-ical-syncer/internal/model"
	"github.com/ronazst/notion-ical-syncer/internal/util"
	"github.com/sirupsen/logrus"
)

func QueryNotionConfigs(tableName string, configIds []string) ([]model.NotionConfig, error) {
	var results []model.NotionConfig
	ddbClient := dynamodb.New(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})))

	for _, configId := range configIds {
		logger := logrus.WithField("config_id", configId)
		logger.Info("Start to query notion config with config id")

		getItemOutput, err := ddbClient.GetItem(&dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key:       map[string]*dynamodb.AttributeValue{"config_id": {S: aws.String(configId)}},
		})
		if err != nil {
			logger.WithError(err).Error("Failed to query dynamodb")
			return nil, util.NewInternalError(fmt.Sprintf("Failed to query dynamodb with id: %s", configId))
		}
		if getItemOutput.Item == nil {
			return nil, util.NewUserInputError(fmt.Sprintf("Notion config not found with id: %s", configId))
		}
		notionConfig := model.NotionConfig{}
		err = dynamodbattribute.UnmarshalMap(getItemOutput.Item, &notionConfig)
		if err != nil {
			logger.WithError(err).Error("Failed unmarshal dynamodb item")
			return nil, util.NewInternalError("Failed to mapping dynamodb item to config struct")
		}
		err = notionConfig.Validate()
		if err != nil {
			logger.WithField("config", fmt.Sprintf("%+v", notionConfig)).Error("AAAAAAAAAAAAAAAAAAAAA")
			logger.WithError(err).Error("Invalid notion config")
			return nil, util.NewInternalError("Invalid notion config")
		}
		results = append(results, notionConfig)
	}

	return results, nil
}

func PutNotionConfig(tableName string, config model.NotionConfig) error {
	logrus.WithField("config", fmt.Sprintf("%+v", config)).Error("BBBBBBBBBBBBBBBBBBBBBBB")
	err := config.Validate()
	if err != nil {
		logrus.WithError(err).Error("Invalid notion config")
		return err
	}

	ddbValue, err := dynamodbattribute.MarshalMap(config)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal dynamodb value")
		return util.NewInternalError("Failed to marshal dynamodb value")
	}

	ddbClient := dynamodb.New(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})))

	_, err = ddbClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      ddbValue,
	})
	if err != nil {
		logrus.WithError(err).Error("Failed to put item to dynamodb")
		return util.NewInternalError("Failed to put item to dynamodb")
	}

	return nil
}
