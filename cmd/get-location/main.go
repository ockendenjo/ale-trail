package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

func (h *LambdaHandler) handle(ctx *handler.Context, event events.APIGatewayProxyRequest) (*data.Details, error) {
	date := event.PathParameters["date"]
	if date == "" {
		return nil, apighandler.HttpError{StatusCode: http.StatusBadRequest, Message: "invalid request", CWLogMsg: "date not set"}
	}
	key := event.QueryStringParameters["key"]
	if key == "" {
		return nil, apighandler.HttpError{StatusCode: http.StatusBadRequest, Message: "invalid request", CWLogMsg: "key not set"}
	}

	foundStash, err := h.getStashForDate(ctx, date, key)
	if err != nil {
		return nil, err
	}

	var details *data.Details
	if foundStash.Details != nil {
		details = foundStash.Details
	} else {
		details = &data.Details{}
	}

	return details, nil
}

func (h *LambdaHandler) getStashForDate(ctx *handler.Context, date, key string) (*data.Stash, error) {
	logger := ctx.GetLogger()
	fileWithEtag, err := h.Client.GetDataFile(ctx)
	if err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusInternalServerError, CWLogErr: err}
	}
	logger.Info("Loaded data file from S3")

	var foundStash *data.Stash
	for _, stash := range fileWithEtag.File.Stashes {
		if stash.Date == "date" {
			foundStash = stash
			break
		}
	}

	if foundStash == nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusNotFound, Message: "no stash found", CWLogMsg: "stash not found"}
	}
	logger.Info("Found stash", "date", date)

	switch key {
	case foundStash.StashKey:
		logger.Info("Authorized using stash key")
	case fileWithEtag.File.AdminKey:
		logger.Info("Authorized using admin key")
	default:
		return nil, apighandler.HttpError{StatusCode: http.StatusUnauthorized, Message: "incorrect key", CWLogMsg: "incorrect stash key"}
	}

	return foundStash, nil
}
