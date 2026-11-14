package main

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	sesTypes "github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/ockendenjo/ale-trail/pkg/data"
	"github.com/ockendenjo/handler"
)

func main() {
	handler.BuildAndStart(func(awsConfig aws.Config) handler.Handler[eventDetail, any] {
		alarmTopicARN := handler.MustGetEnv("ALARM_TOPIC_ARN")
		bucketName := handler.MustGetEnv("DATA_BUCKET")
		replyTo := handler.MustGetEnv("REPLY_TO")
		domain := handler.MustGetEnv("DOMAIN")

		sesClient := ses.NewFromConfig(awsConfig)
		s3Client := s3.NewFromConfig(awsConfig)
		snsClient := sns.NewFromConfig(awsConfig)

		h := &lambdaHandler{
			dataClient:    data.NewClient(s3Client, bucketName),
			sesClient:     sesClient,
			replyTo:       replyTo,
			snsClient:     snsClient,
			alarmTopicARN: alarmTopicARN,
			domain:        domain,
		}
		return h.handle
	})
}

type lambdaHandler struct {
	dataClient    data.Client
	sesClient     *ses.Client
	alarmTopicARN string
	snsClient     *sns.Client
	domain        string
	replyTo       string
}

func (h *lambdaHandler) handle(ctx *handler.Context, detail eventDetail) (any, error) {
	ctx.GetLogger().AddParam("input", detail)

	fileWithEtag, err := h.dataClient.GetDataFile(ctx)
	if err != nil {
		return nil, err
	}

	ctx.GetLogger().AddParam("template", detail.Template)

	switch detail.Template {
	case "DAY-2":
		return h.dayMinus2(ctx, fileWithEtag)
	case "DAY-1":
		return h.dayMinus1(ctx, fileWithEtag)
	case "OVERDUE":
		return h.overdue(ctx, fileWithEtag)
	default:
		return nil, errors.New("unknown email template")
	}
}

func (h *lambdaHandler) sendTopicMessage(ctx *handler.Context, message string) error {
	_, err := h.snsClient.Publish(ctx, &sns.PublishInput{
		TopicArn: &h.alarmTopicARN,
		Message:  &message,
	})
	if err != nil {
		return err
	}

	ctx.GetLogger().AddParam("snsMsg", message).Info("Sent message to topic")
	return err
}

func (h *lambdaHandler) sendEmail(ctx *handler.Context, toAddr, subject, message string) error {
	_, err := h.sesClient.SendEmail(ctx, &ses.SendEmailInput{
		Destination: &sesTypes.Destination{
			ToAddresses: []string{toAddr},
		},
		Message: &sesTypes.Message{
			Subject: &sesTypes.Content{Data: &subject},
			Body:    &sesTypes.Body{Text: &sesTypes.Content{Data: &message}},
		},
		Source:           new(fmt.Sprintf("noreply@%s", h.domain)),
		ReplyToAddresses: []string{h.replyTo},
	})
	if err != nil {
		return err
	}

	ctx.GetLogger().Info("Sent email")
	return nil
}

type eventDetail struct {
	Template string `json:"template"`
}
