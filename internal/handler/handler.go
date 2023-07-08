package handler

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ronazst/notion-ical-syncer/internal/awsutil"
	"github.com/ronazst/notion-ical-syncer/internal/ical"
	"github.com/ronazst/notion-ical-syncer/internal/model"
	"github.com/ronazst/notion-ical-syncer/internal/notion"
	"github.com/ronazst/notion-ical-syncer/internal/util"
	"github.com/sirupsen/logrus"
	"strings"
)

var DefaultInternalError = util.NewInternalError(util.MessageInternalServerError)

func handleRequest(ctx context.Context, event events.LambdaFunctionURLRequest) (string, error) {
	stackId, err := util.GetOsEnv(util.EnvStackId)
	if err != nil {
		logrus.WithError(err).Error("Failed to get stack id")
		return "", DefaultInternalError
	}
	ddbTableName, err := util.GetOsEnv(util.EnvDdbTable)
	if err != nil {
		logrus.WithError(err).Error("Failed to get ddb table name")
		return "", DefaultInternalError
	}

	if err := util.IsEligibleUser(event.RawPath, stackId); err != nil {
		logrus.WithError(err).WithField("raw_path", event.RawPath).Warn("Request not start with stack id")
		return "", DefaultInternalError
	}

	configIds := event.QueryStringParameters[util.PathArgConfigIds]
	if util.IsBlank(configIds) {
		logrus.WithField("stack_id", stackId).Warn("User request without config_ids")
		return "", util.NewUserInputError("query parameter without config_ids")
	}

	calendarData := make([]model.CalendarData, 0)
	for _, configId := range strings.Split(configIds, ",") {
		logger := logrus.WithField("config_id", configId)

		logger.Info("Start to query notion config with config id")
		notionConfig, err := awsutil.QueryNotionConfig(ddbTableName, configId)
		if err != nil {
			logger.WithError(err).Error("Failed to query notion config")
			continue
		}

		logger = logger.WithField("notion_db_id", notionConfig.NotionDbId)

		logger.Info("Start to query calendar data with notion config")
		data, err := notion.QueryCalendarData(ctx, notionConfig.NotionToken, []model.QueryItem{
			{DatabaseID: notionConfig.NotionDbId, DateFieldKey: notionConfig.FieldMapping.Date},
		})
		if err != nil {
			logger.WithError(err).Error("Failed to query calendar data with notion config")
			return "", err
		}
		logger.Info("Successfully query calendar data")
		calendarData = append(calendarData, data...)
	}

	calContent, err := ical.TransToICal(calendarData)
	if err != nil {
		logrus.WithError(err).Error("Failed to convert calendar data to ICal format")
		return "", err
	}

	return calContent, nil
}
