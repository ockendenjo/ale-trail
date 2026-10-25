package teststubs

import (
	"context"

	"github.com/ockendenjo/ale-trail/pkg/data"
)

type MockDataClient struct {
	GetDataFileFn       func(ctx context.Context) (*data.FileWithETag, error)
	PutFileWithEtagData *data.FileWithETag
	PutFileWithEtagFn   func(ctx context.Context, fileWithEtag *data.FileWithETag) error
}

func (m *MockDataClient) GetDataFile(ctx context.Context) (*data.FileWithETag, error) {
	return m.GetDataFileFn(ctx)
}

func (m *MockDataClient) PutFileWithEtag(ctx context.Context, fileWithETag *data.FileWithETag) error {
	m.PutFileWithEtagData = fileWithETag
	return m.PutFileWithEtagFn(ctx, fileWithETag)
}
