package data

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client interface {
	GetDataFile(ctx context.Context) (*FileWithETag, error)
}

func NewClient(s3Client *s3.Client, bucketName string) Client {
	return &client{
		s3Client:   s3Client,
		bucketName: bucketName,
	}
}

type client struct {
	s3Client   *s3.Client
	bucketName string
}

func (c *client) GetDataFile(ctx context.Context) (*FileWithETag, error) {

	res, err := c.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &c.bucketName,
		Key:    new("data.json"),
	})
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	var file *File
	err = json.NewDecoder(res.Body).Decode(file)
	if err != nil {
		return nil, err
	}

	return &FileWithETag{
		File: file,
		ETag: res.ETag,
	}, nil
}
