package web

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Client interface {
	FileExists(ctx context.Context, path string) (bool, error)
}

func NewClient(s3Client *s3.Client, bucketName string) Client {
	return &client{s3Client: s3Client, bucketName: bucketName}
}

type client struct {
	s3Client   *s3.Client
	bucketName string
}

func (c *client) FileExists(ctx context.Context, path string) (bool, error) {

	_, err := c.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &c.bucketName,
		Key:    &path,
	})

	if err != nil {
		if _, b := errors.AsType[*s3Types.NoSuchKey](err); b {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
