package e2e

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/failsafe-go/failsafe-go"
	"github.com/failsafe-go/failsafe-go/retrypolicy"
	"github.com/ockendenjo/ale-trail/pkg/data"
	"github.com/ockendenjo/handler"
	g "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
)

type LogHelper struct {
	t          g.FullGinkgoTInterface
	dataClient data.Client
}

func NewLogHelper(t g.FullGinkgoTInterface) *LogHelper {
	awsConfig, err := config.LoadDefaultConfig(t.Context())
	assert.NoError(t, err)

	s3Client := s3.NewFromConfig(awsConfig)
	dataClient := data.NewClient(s3Client, handler.MustGetEnv("DATA_BUCKET"))

	return &LogHelper{
		t:          t,
		dataClient: dataClient,
	}
}

func (h *LogHelper) PutDataFile(f *data.File) {
	fileWithEtag := data.FileWithETag{File: f}
	err := h.dataClient.PutFileWithEtag(h.t.Context(), &fileWithEtag)
	assert.NoError(h.t, err)
}

func (h *LogHelper) GetDataFile() *data.File {
	fileWithETag, err := h.dataClient.GetDataFile(h.t.Context())
	assert.NoError(h.t, err)
	return fileWithETag.File
}

func (h *LogHelper) AwaitDataFile(p func(*data.File) bool) *data.File {

	getFile := func() (*data.File, error) {
		fileWithETag, err := h.dataClient.GetDataFile(h.t.Context())
		if err != nil {
			return nil, err
		}
		return fileWithETag.File, nil
	}

	policy := retryWithPredicate(p)

	file, err := failsafe.With(policy).Get(getFile)

	assert.NoError(h.t, err)
	return file
}

func retryWithPredicate[T any](p func(t T) bool) retrypolicy.RetryPolicy[T] {
	return retrypolicy.NewBuilder[T]().
		HandleIf(func(data T, err error) bool {
			if err != nil {
				return true
			}
			return !p(data)
		}).
		WithMaxRetries(10).
		WithBackoff(100*time.Millisecond, time.Second).
		Build()
}
