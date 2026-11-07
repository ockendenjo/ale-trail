package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ockendenjo/handler"
)

func main() {
	handler.BuildAndStart(func(awsConfig aws.Config) handler.Handler[events.APIGatewayProxyRequest, events.APIGatewayProxyResponse] {
		h := &lambdaHandler{
			s3Client: s3.NewFromConfig(awsConfig),
			bucket:   handler.MustGetEnv("WEB_BUCKET"),
		}
		return h.handle
	})
}

type lambdaHandler struct {
	s3Client *s3.Client
	bucket   string
}

func (h *lambdaHandler) handle(ctx *handler.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	obj, err := h.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(h.bucket),
		Key:    aws.String("index.html"),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, fmt.Errorf("fetching index.html: %w", err)
	}
	defer obj.Body.Close()

	htmlBytes, err := io.ReadAll(obj.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, fmt.Errorf("reading index.html: %w", err)
	}

	nonceBytes := make([]byte, 16)
	if _, err = rand.Read(nonceBytes); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, fmt.Errorf("generating nonce: %w", err)
	}
	nonce := base64.StdEncoding.EncodeToString(nonceBytes)

	html := strings.ReplaceAll(string(htmlBytes), "__NONCE__", nonce)

	csp := strings.Join([]string{
		"default-src 'self'",
		"script-src 'self' 'nonce-" + nonce + "'",
		"style-src 'self' 'unsafe-inline'",
		"font-src 'self'",
		"img-src 'self' data:",
		"frame-ancestors 'none'",
	}, "; ")

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":            "text/html; charset=utf-8",
			"Cache-Control":           "no-store",
			"Content-Security-Policy": csp,
		},
		Body: html,
	}, nil
}
