package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ockendenjo/ale-trail/pkg/api"
	"github.com/ockendenjo/ale-trail/pkg/apighandler"
	"github.com/ockendenjo/ale-trail/pkg/data"
	"github.com/ockendenjo/handler"
)

func main() {

	handler.BuildAndStart(func(awsConfig aws.Config) apighandler.H {
		bucketName := handler.MustGetEnv("DATA_BUCKET")
		s3Client := s3.NewFromConfig(awsConfig)

		h := &LambdaHandler{
			Client: data.NewClient(s3Client, bucketName),
		}
		return apighandler.GetHandler(h.handle, http.StatusOK)
	})
}

type LambdaHandler struct {
	Client data.Client
}

func (h *LambdaHandler) handle(ctx *handler.Context, event events.APIGatewayProxyRequest) (*api.Stash, error) {
	logger := ctx.GetLogger()

	date := event.PathParameters["date"]
	if date == "" {
		return nil, apighandler.HttpError{StatusCode: http.StatusBadRequest, Message: "invalid request", CWLogMsg: "date not set"}
	}
	logger.AddParam("date", date)

	fileWithEtag, err := h.Client.GetDataFile(ctx)
	if err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusInternalServerError, CWLogErr: err}
	}
	logger.Info("Loaded data file from S3")

	foundStash := fileWithEtag.GetStashForDate(date)
	if foundStash == nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusNotFound, Message: "not found", CWLogMsg: "not found"}
	}

	stash := api.GetStash(foundStash, true)
	return stash, nil
}
