package handler

import (
	"context"
	ics "github.com/arran4/golang-ical"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ronazst/notion-ical-syncer/internal/awsutil"
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

	configIdsStr := strings.TrimSpace(event.QueryStringParameters[util.PathArgConfigIds])
	if len(configIdsStr) == 0 {
		logrus.WithField("stack_id", stackId).Warn("User request without config_ids")
		return "", util.NewUserInputError("query parameter without config_ids")
	}

	notionConfigs, err := awsutil.QueryNotionConfigs(ddbTableName, strings.Split(configIdsStr, ","))
	if err != nil {
		return "", err
	}

	vEvents, err := notion.QueryEvents(ctx, notionConfigs)
	if err != nil {
		return "", err
	}

	cal := ics.NewCalendar()
	cal.SetProductId("notion-ical-syncer")
	cal.SetVersion("2.0")
	for _, vEvent := range vEvents {
		cal.AddVEvent(vEvent)
	}

	return cal.Serialize(), nil
}
