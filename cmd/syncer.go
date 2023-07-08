package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ronazst/notion-ical-syncer/internal/handler"
	"github.com/ronazst/notion-ical-syncer/internal/log"
)

func main() {
	lambda.Start(func(ctx context.Context, event events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
		log.SetupLogrus(event.RequestContext.RequestID)
		return handler.Handle(ctx, event)
	})
}
