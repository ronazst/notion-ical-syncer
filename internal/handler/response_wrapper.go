package handler

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ronazst/notion-ical-syncer/internal/util"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Handle(ctx context.Context, event events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	logrus.WithField("event", fmt.Sprintf("%+v", event)).Info("Query ICal with event")
	responseContent, err := handleRequest(ctx, event)
	response := events.LambdaFunctionURLResponse{
		StatusCode: http.StatusOK,
		Body:       responseContent,
	}
	logrus.WithError(err).Info("Finished handling request")

	if err != nil {
		switch value := err.(type) {
		case *util.Error:
			response.StatusCode = value.Code()
			response.Body = value.Error()
		default:
			response.StatusCode = http.StatusInternalServerError
			response.Body = util.MessageInternalServerError
		}
	}

	return response, nil
}
