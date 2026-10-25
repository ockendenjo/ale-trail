package main

import (
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ockendenjo/ale-trail/pkg/apighandler"
	"github.com/ockendenjo/ale-trail/pkg/card"
	"github.com/ockendenjo/ale-trail/pkg/data"
	"github.com/ockendenjo/handler"
)

func main() {
	handler.BuildAndStart(func(awsConfig aws.Config) apighandler.H {
		bucketName := handler.MustGetEnv("DATA_BUCKET")
		s3Client := s3.NewFromConfig(awsConfig)

		h := &lambdaHandler{
			Client: data.NewClient(s3Client, bucketName),
		}
		return apighandler.GetHandler(h.handle, http.StatusOK)
	})
}

type lambdaHandler struct {
	Client     data.Client
	BucketName string
	S3Client   *s3.Client
}

func (h *lambdaHandler) handle(ctx *handler.Context, _ events.APIGatewayProxyRequest) (*CardResponse, error) {
	logger := ctx.GetLogger()

	res, err := h.Client.GetDataFile(ctx)
	if err != nil {
		return nil, err
	}
	logger.Info("Fetched data file from S3")

	cardRes := getResponse(res.File)
	logger.Info("Mapped data to cards")

	return cardRes, nil
}

func getResponse(file *data.File) *CardResponse {
	a := make([]card.Card, 0, len(file.Stashes))

	for _, stash := range file.Stashes {
		c := stashToCard(stash)
		a = append(a, c)
	}

	return &CardResponse{Cards: a}
}

func stashToCard(stash *data.Stash) card.Card {
	locked := isLocked(stash)
	var image *string
	if !locked {
		image = &stash.Image
	}
	var releaseTime *string
	if locked && stash.ReleaseType == "TIMED" && isToday(stash) {
		releaseTime = &stash.ReleaseTime
	}

	c := card.Card{
		Date:        stash.Date,
		Locked:      locked,
		ReleaseType: new(stash.ReleaseType),
		ReleaseTime: releaseTime,
		Image:       image,
		HostNeeded:  isHostNeeded(stash),
	}
	return c
}

func isToday(stash *data.Stash) bool {
	return stash.Date == time.Now().Format("2006-01-02")
}

func isLocked(stash *data.Stash) bool {
	s := stash.Date + "T" + stash.ReleaseTime + ":00Z"
	tm, err := time.Parse("2006-01-02T15:04:05Z", s)
	if err != nil {
		return true
	}
	return tm.After(time.Now())
}

func isHostNeeded(stash *data.Stash) *bool {
	if stash.Stasher != nil {
		return nil
	}
	return new(true)
}

type CardResponse struct {
	Cards []card.Card `json:"cards"`
}
