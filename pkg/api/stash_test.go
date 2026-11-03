package api

import (
	"testing"

	"github.com/ockendenjo/ale-trail/pkg/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetStash(t *testing.T) {

	testcases := []struct {
		name        string
		fullMode    bool
		input       *data.Stash
		checkResult func(t *testing.T, stash *Stash)
	}{
		{
			name:     "should include image in short mode",
			fullMode: false,
			input: &data.Stash{
				Image:       "sample.jpeg",
				Date:        "2026-11-02",
				ReleaseTime: "12:00",
				Details:     new(data.Details),
			},
			checkResult: func(t *testing.T, stash *Stash) {
				require.NotNil(t, stash.Image)
				assert.Equal(t, "sample.jpeg", *stash.Image)
			},
		},
		{
			name: "should exclude image if stash is locked",
			fullMode: false,
			input: &data.Stash{
				Image:       "sample.jpeg",
				Date:        "2026-11-02",
				ReleaseTime: "12:00",
				Details:     nil,
			},
			checkResult: func(t *testing.T, stash *Stash) {
				require.Nil(t, stash.Image)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			res := GetStash(tc.input, tc.fullMode)
			tc.checkResult(t, res)
		})
	}
}
