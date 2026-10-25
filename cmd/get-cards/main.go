package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ockendenjo/ale-trail/pkg/apighandler"
	"github.com/ockendenjo/ale-trail/pkg/card"
	"github.com/ockendenjo/handler"
)

func main() {
	handler.BuildAndStart(func(awsConfig aws.Config) apighandler.H {
		bucketName := handler.MustGetEnv("DATA_BUCKET")
		s3Client := s3.NewFromConfig(awsConfig)

		h := &lambdaHandler{
			Client: card.NewClient(s3Client, bucketName),
		}
		return apighandler.GetHandler(h.handle, http.StatusOK)
	})
}

type lambdaHandler struct {
	Client     card.Client
	BucketName string
	S3Client   *s3.Client
}

func (h *lambdaHandler) handle(ctx *handler.Context, _ events.APIGatewayProxyRequest) (*CardResponse, error) {

	res, err := h.Client.GetCards(ctx)
	if err != nil {
		return nil, err
	}

	return &CardResponse{Cards: res}, nil
}

type CardResponse struct {
	Cards []card.Card `json:"cards"`
}
