package main

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/ockendenjo/ale-trail/pkg/apighandler"
	"github.com/ockendenjo/ale-trail/pkg/data"
	"github.com/ockendenjo/ale-trail/pkg/testing/teststubs"
	"github.com/ockendenjo/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_handleRequest(t *testing.T) {

	testcases := []struct {
		name        string
		setupClient func(c *teststubs.MockDataClient)
		checkFn     func(t *testing.T, c *teststubs.MockDataClient, err error)
	}{
		{
			name: "success",
			setupClient: AllOf(
				MockDataOf([]*data.Stash{
					{
						Date: "2026-10-24",
					},
					{
						Date:     "2026-10-25",
						StashKey: "oldKey",
					},
				}),
				PutOk,
			),
			checkFn: func(t *testing.T, c *teststubs.MockDataClient, err error) {
				require.NoError(t, err)

				assert.NotNil(t, c.PutFileWithEtagData)
				s := c.PutFileWithEtagData.File.Stashes[1]

				assert.NotNil(t, s.Stasher)
				assert.Equal(t, "John", s.Stasher.Name)
				assert.Equal(t, "john@example.com", s.Stasher.Email)
				assert.NotEqual(t, "oldKey", s.StashKey)
			},
		},
		{
			name: "stash not found",
			setupClient: MockDataOf([]*data.Stash{
				{
					Date: "2026-10-24",
				},
			}),
			checkFn: IsHttpError(http.StatusNotFound),
		},
		{
			name: "stash already has stasher",
			setupClient: MockDataOf([]*data.Stash{
				{
					Date: "2026-10-25",
					Stasher: &data.Stasher{
						Name:  "Anne",
						Email: "anne@example.com",
					},
				},
			}),
			checkFn: IsHttpError(http.StatusConflict),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			dataClient := &teststubs.MockDataClient{}
			tc.setupClient(dataClient)

			h := &lambdaHandler{Client: dataClient}
			_, err := h.handleRequest(handler.Get(t.Context()), "2026-10-25", &SetStasherRequest{Name: "John", Email: "john@example.com"})
			tc.checkFn(t, dataClient, err)
		})
	}
}

func IsHttpError(statusCode int) func(t *testing.T, _ *teststubs.MockDataClient, err error) {
	return func(t *testing.T, _ *teststubs.MockDataClient, err error) {
		httpErr, isErr := errors.AsType[apighandler.HttpError](err)
		require.True(t, isErr, "error not expected type")
		assert.Equal(t, statusCode, httpErr.StatusCode)
	}
}

func AllOf(f ...func(*teststubs.MockDataClient)) func(*teststubs.MockDataClient) {
	return func(dataClient *teststubs.MockDataClient) {
		for _, f := range f {
			f(dataClient)
		}
	}
}

func MockDataOf(stashes []*data.Stash) func(*teststubs.MockDataClient) {
	return func(c *teststubs.MockDataClient) {
		c.GetDataFileFn = func(ctx context.Context) (*data.FileWithETag, error) {
			return &data.FileWithETag{
				File: &data.File{
					Schema:   nil,
					AdminKey: "",
					Stashes:  stashes,
				},
				ETag: new("etag"),
			}, nil
		}
	}
}

func PutOk(dataClient *teststubs.MockDataClient) {
	dataClient.PutFileWithEtagFn = func(ctx context.Context, fileWithETag *data.FileWithETag) error {
		return nil
	}
}
