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
	"github.com/ockendenjo/ale-trail/pkg/web"
	"github.com/ockendenjo/handler"
)

func main() {
	handler.BuildAndStart(func(awsConfig aws.Config) handler.Handler[eventDetail, any] {
		alarmTopicARN := handler.MustGetEnv("ALARM_TOPIC_ARN")
		dataBucket := handler.MustGetEnv("DATA_BUCKET")
		webBucket := handler.MustGetEnv("WEB_BUCKET")
		replyTo := handler.MustGetEnv("REPLY_TO")
		domain := handler.MustGetEnv("DOMAIN")
		awsEnv := handler.MustGetEnv("AWS_ENV")

		sesClient := ses.NewFromConfig(awsConfig)
		s3Client := s3.NewFromConfig(awsConfig)
		snsClient := sns.NewFromConfig(awsConfig)

		h := &lambdaHandler{
			dataClient:    data.NewClient(s3Client, dataBucket),
			webClient:     web.NewClient(s3Client, webBucket),
			sesClient:     sesClient,
			replyTo:       replyTo,
			snsClient:     snsClient,
			alarmTopicARN: alarmTopicARN,
			domain:        domain,
			awsEnv:        awsEnv,
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
	awsEnv        string
	webClient     web.Client
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
	case "IMAGE-CHECK":
		return h.imageCheck(ctx, fileWithEtag)
	case "TIDY-UP":
		return h.tidyUp(ctx, fileWithEtag)
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
		Source:           h.getSource(),
		ReplyToAddresses: []string{h.replyTo},
	})
	if err != nil {
		return err
	}

	ctx.GetLogger().Info("Sent email")
	return nil
}

func (h *lambdaHandler) getSource() *string {
	prefix := ""
	if h.awsEnv == "dev" {
		prefix = "[dev] "
	}

	source := fmt.Sprintf("%sHBT Ale Trail <noreply@%s>", prefix, h.domain)
	return &source
}

type eventDetail struct {
	Template string `json:"template"`
}
