package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdaurl"
	"github.com/ronazst/notion-ical-syncer/internal/util"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/ronazst/notion-ical-syncer/internal/handler"
)

func main() {
	lambda.Start(func(ctx context.Context, event *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLStreamingResponse, error) {
		util.SetupLogrus(event.RequestContext.RequestID)
		iCalHandler, err := handler.NewHandler()
		if err != nil {
			logrus.WithError(err).Error("failed to create handler")
			return &events.LambdaFunctionURLStreamingResponse{
				StatusCode: http.StatusInternalServerError,
			}, nil
		}
		return lambdaurl.Wrap(iCalHandler)(ctx, event)
	})
}
