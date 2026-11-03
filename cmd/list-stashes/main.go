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

func (h *LambdaHandler) handle(ctx *handler.Context, _ events.APIGatewayProxyRequest) (*api.Stashes, error) {
	logger := ctx.GetLogger()

	fileWithEtag, err := h.Client.GetDataFile(ctx)
	if err != nil {
		return nil, apighandler.HttpError{StatusCode: http.StatusInternalServerError, CWLogErr: err}
	}
	logger.Info("Loaded data file from S3")

	a := make([]*api.Stash, 0, len(fileWithEtag.File.Stashes))
	for _, s := range fileWithEtag.File.Stashes {
		a = append(a, api.GetStash(s, false))
	}

	return &api.Stashes{Stashes: a}, nil
}
