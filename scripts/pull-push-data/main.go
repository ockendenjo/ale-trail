package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ockendenjo/ale-trail/pkg/data"
	"github.com/ockendenjo/handler"
)

func main() {
	env := getEnv()

	ctx := context.Background()
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	s3Client := s3.NewFromConfig(awsConfig)
	bucketName := getBucketName(env)
	dataClient := data.NewClient(s3Client, bucketName)

	fileWithEtag, err := dataClient.GetDataFile(ctx)
	if err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(fileWithEtag.File, "", "    ")
	if err != nil {
		panic(err)
	}

	fileName := getLocalFileName(env)
	if err = os.WriteFile(fileName, b, 0o644); err != nil { // #nosec G306
		panic(err)
	}

	fmt.Printf("Data file written to %s\n", fileName)
	fmt.Println("Press enter when editing is complete")
	_, _ = fmt.Scanln()

	b, err = os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	var file *data.File
	if err = json.Unmarshal(b, &file); err != nil {
		panic(err)
	}

	fileWithEtag.File = file

	if err = dataClient.PutFileWithEtag(ctx, fileWithEtag); err != nil {
		panic(err)
	}

	fmt.Println("Data file uploaded to S3")
	fmt.Println("Success :)")
}

func getLocalFileName(env string) string {
	return fmt.Sprintf("%s.data.json", env)
}

func getBucketName(env string) string {
	if env == "dev" {
		return handler.MustGetEnv("DEV_DATA_BUCKET")
	}
	return handler.MustGetEnv("PRO_DATA_BUCKET")
}

func getEnv() string {
	if len(os.Args) < 2 {
		panic("missing env arg")
	}
	env := os.Args[1]
	if env == "dev" || env == "pro" {
		return env
	}
	panic("invalid env value")
}
