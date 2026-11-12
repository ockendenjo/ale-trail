package main

import (
	"net/http"
	"uuid"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ockendenjo/ale-trail/pkg/apighandler"
	"github.com/ockendenjo/ale-trail/pkg/data"
	"github.com/ockendenjo/ale-trail/pkg/util"
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

	fileWithEtag, err := h.Client.GetDataFile(ctx)
	if err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusInternalServerError, CWLogErr: err}
	}
	logger.Info("Loaded data file from S3")

	foundStash, err := util.GetStashForDate(logger, fileWithEtag, date, key)
	if err != nil {
		return nil, err
	}

	if foundStash.Details != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusBadRequest, Message: "cannot remove stasher when a location is set", CWLogMsg: "stash had a location set - refusing to delete stasher"}
	}

	foundStash.Stasher = nil
	foundStash.StashKey = uuid.New().String()
	logger.Info("Removed stasher and rotated stash key")

	err = h.Client.PutFileWithEtag(ctx, fileWithEtag)
	if err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusInternalServerError, CWLogErr: err, CWLogMsg: "failed to write data file to S3"}
	}
	logger.Info("Wrote data file to S3")
	return nil, nil
}
