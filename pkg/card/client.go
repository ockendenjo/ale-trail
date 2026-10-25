package card

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client interface {
	GetCards(ctx context.Context) ([]Card, error)
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

func (c *client) GetCards(ctx context.Context) ([]Card, error) {

	res, err := c.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &c.bucketName,
		Key:    new("cards.json"),
	})
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	var cards []Card
	err = json.NewDecoder(res.Body).Decode(&cards)
	if err != nil {
		return nil, err
	}

	return cards, nil
}
