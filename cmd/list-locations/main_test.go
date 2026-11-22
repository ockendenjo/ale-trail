package main

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ockendenjo/ale-trail/pkg/data"
	"github.com/ockendenjo/ale-trail/pkg/testing/teststubs"
	"github.com/ockendenjo/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_handler(t *testing.T) {

	now := time.Now()

	getOffsetDate := func(n int) string {
		return now.AddDate(0, 0, n).Format(time.DateOnly)
	}

	stashes := []*data.Stash{
		{Date: getOffsetDate(-2), ReleaseTime: "12:00", Details: &data.Details{Lat: 56.999, Lon: -2.121}, Image: "image.jpeg"}, //Details present = unlocked
		{Date: getOffsetDate(-1), ReleaseTime: "12:00", Details: &data.Details{Lat: 56.123, Lon: -3.192}},                      //Missing image = locked
		{Date: getOffsetDate(1), ReleaseTime: "12:00", Details: &data.Details{}, Image: "image.jpeg"},                          //In the future = locked
	}

	assert.False(t, stashes[0].IsLocked(now))
	assert.True(t, stashes[1].IsLocked(now))
	assert.True(t, stashes[2].IsLocked(now))

	mockDataClient := &teststubs.MockDataClient{
		GetDataFileFn: func(ctx context.Context) (*data.FileWithETag, error) {
			return &data.FileWithETag{File: &data.File{Stashes: stashes}}, nil
		},
	}

	h := LambdaHandler{Client: mockDataClient}

	res, err := h.handle(handler.Get(t.Context()), events.APIGatewayProxyRequest{})
	require.NoError(t, err)

	exp := Response{
		Locations: [][2]float64{
			{-2.121, 56.999},
		},
	}
	assert.Equal(t, exp, *res)
}
