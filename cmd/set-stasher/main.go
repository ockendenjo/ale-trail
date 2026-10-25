package main

import (
	"encoding/json/v2"
	"errors"
	"net/http"
	"uuid"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/ockendenjo/ale-trail/pkg/apighandler"
	"github.com/ockendenjo/ale-trail/pkg/data"
	"github.com/ockendenjo/handler"
)

func main() {
	handler.BuildAndStart(func(awsConfig aws.Config) apighandler.H {
		bucketName := handler.MustGetEnv("DATA_BUCKET")
		s3Client := s3.NewFromConfig(awsConfig)

		h := &lambdaHandler{Client: data.NewClient(s3Client, bucketName)}
		return apighandler.GetHandler(h.handle, http.StatusNoContent)
	})
}

type lambdaHandler struct {
	Client data.Client
}

func (h *lambdaHandler) handle(ctx *handler.Context, event events.APIGatewayProxyRequest) (any, error) {
	var req *SetStasherRequest
	err := json.Unmarshal([]byte(event.Body), &req)
	if err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusBadRequest, Message: "invalid request"}
	}
	if err := req.Validate(); err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusBadRequest, Message: err.Error(), CWLogMsg: "request failed validation"}
	}
	date := event.PathParameters["date"]
	if date == "" {
		return nil, apighandler.HttpError{StatusCode: http.StatusBadRequest, Message: "invalid request", CWLogMsg: "date not set"}
	}

	return h.handleRequest(ctx, date, req)
}

func (h *lambdaHandler) handleRequest(ctx *handler.Context, date string, req *SetStasherRequest) (any, error) {
	logger := ctx.GetLogger()

	fileWithEtag, err := h.Client.GetDataFile(ctx)
	if err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusInternalServerError, CWLogErr: err}
	}
	logger.Info("Loaded data file from S3")

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
	logger.Info("Found stash", "date", date)
	if foundStash.Stasher != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusConflict, Message: "stash has stasher already", CWLogMsg: "stash already had stasher"}
	}

	//Update fields
	foundStash.Stasher = &data.Stasher{
		Name:  req.Name,
		Email: req.Email,
	}
	logger.Info("Set stasher info")

	foundStash.StashKey = uuid.New().String()
	logger.Info("Updated stash key")

	//Write file back
	err = h.Client.PutFileWithEtag(ctx, fileWithEtag)
	if err != nil {
		if isPreconditionFailed(err) {
			return nil, apighandler.HttpError{StatusCode: http.StatusPreconditionFailed, CWLogMsg: "file in S3 changed since read", CWLogErr: err}
		}
		return nil, apighandler.HttpError{StatusCode: http.StatusInternalServerError, CWLogErr: err}
	}
	logger.Info("Wrote data file to S3")
	return nil, nil
}

func isPreconditionFailed(err error) bool {
	var respErr *smithyhttp.ResponseError
	return errors.As(err, &respErr) && respErr.HTTPStatusCode() == http.StatusPreconditionFailed
}
