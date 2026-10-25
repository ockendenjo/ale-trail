package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/ockendenjo/ale-trail/pkg/apighandler"
	"github.com/ockendenjo/ale-trail/pkg/card"
	"github.com/ockendenjo/handler"
)

func main() {
	handler.BuildAndStart(func(awsConfig aws.Config) apighandler.H {

		h := &lambdaHandler{}
		return apighandler.GetHandler(h.handle, http.StatusOK)
	})
}

type lambdaHandler struct{}

func (h *lambdaHandler) handle(*handler.Context, events.APIGatewayProxyRequest) (CardResponse, error) {

	return CardResponse{Cards: make([]card.Card, 0)}, nil
}

type CardResponse struct {
	Cards []card.Card `json:"cards"`
}
