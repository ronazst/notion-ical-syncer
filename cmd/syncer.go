package main

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/notion-ical-syncer/internal/ical"
	"github.com/notion-ical-syncer/internal/model"
	"github.com/notion-ical-syncer/internal/notion"
	"github.com/notion-ical-syncer/internal/util"
	"os"
)

func HandleRequest(ctx context.Context, event events.LambdaFunctionURLRequest) (string, error) {
	token := event.QueryStringParameters["notion_token"]
	queryItems, err := model.ParseQueryItems(event.QueryStringParameters["db_date_pair"])
	if err != nil {
		return "failed to parse query items", err
	}

	if util.IsBlank(token) {
		token = os.Getenv("NOTION_TOKEN")
	}
	if util.IsBlank(token) || len(queryItems) == 0 {
		return "", errors.New("invalid token or db_date_pair")
	}

	calendarData, err := notion.QueryCalendarData(ctx, token, queryItems)
	if err != nil {
		return "failed to query calendar data", err
	}
	cal, err := ical.TransToICal(calendarData)
	if err != nil {
		return "failed to trans to ical", err
	}

	return cal, nil
}

func main() {
	lambda.Start(HandleRequest)
}
