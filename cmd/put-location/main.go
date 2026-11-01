package main

import (
	"encoding/json/v2"
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
		return apighandler.GetHandler(h.handle, http.StatusNoContent)
	})
}

type LambdaHandler struct {
	Client data.Client
}

func (h *LambdaHandler) handle(ctx *handler.Context, event events.APIGatewayProxyRequest) (any, error) {
	logger := ctx.GetLogger()

	date := event.PathParameters["date"]
	if date == "" {
		return nil, apighandler.HttpError{StatusCode: http.StatusBadRequest, Message: "invalid request", CWLogMsg: "date not set"}
	}
	logger.AddParam("date", date)
	key := event.QueryStringParameters["key"]
	if key == "" {
		return nil, apighandler.HttpError{StatusCode: http.StatusBadRequest, Message: "invalid request", CWLogMsg: "key not set"}
	}

	var details *data.Details
	err := json.Unmarshal([]byte(event.Body), &details)
	if err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusBadRequest, Message: "invalid request", CWLogMsg: "JSON conversion failed"}
	}
	if err := details.Validate(); err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusBadRequest, Message: err.Error(), CWLogMsg: "request validation failed"}
	}

	fileWithEtag, err := h.Client.GetDataFile(ctx)
	if err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusInternalServerError, CWLogErr: err}
	}
	logger.Info("Loaded data file from S3")

	foundStash, err := h.getStashForDate(logger, fileWithEtag, date, key)
	if err != nil {
		return nil, err
	}

	foundStash.Details = details

	err = h.Client.PutFileWithEtag(ctx, fileWithEtag)
	if err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusInternalServerError, CWLogErr: err, CWLogMsg: "failed to write data file to S3"}
	}
	logger.Info("Wrote data file to S3")
	return nil, nil
}

func (h *LambdaHandler) getStashForDate(logger *handler.Logger, fileWithEtag *data.FileWithETag, date, key string) (*data.Stash, error) {
	var foundStash *data.Stash
	for _, stash := range fileWithEtag.File.Stashes {
		if stash.Date == date {
			foundStash = stash
			break
		}
	}

	if foundStash == nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusNotFound, Message: "no stash found", CWLogMsg: "stash not found"}
	}
	logger.Info("Found stash")

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
