package awsutil

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/ronazst/notion-ical-syncer/internal/model"
)

func QueryNotionConfig(tableName string, notionConfigId string) (*model.NotionConfig, error) {
	ddbClient := createDdbClient()
	result, err := ddbClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       map[string]*dynamodb.AttributeValue{"config_id": {S: aws.String(notionConfigId)}},
	})
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, errors.New("notion config not found")
	}
	notionConfig := model.NotionConfig{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &notionConfig)
	if err != nil {
		return nil, err
	}

	return &notionConfig, nil
}

func createDdbClient() *dynamodb.DynamoDB {
	return dynamodb.New(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})))
}
