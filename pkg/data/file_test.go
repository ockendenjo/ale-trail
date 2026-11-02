package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getStashForDate(t *testing.T) {
	file := &FileWithETag{
		File: &File{
			Schema:   nil,
			AdminKey: "",
			Stashes: []*Stash{
				{
					Date: "2026-11-02",
				},
			},
		},
		ETag: nil,
	}

	res := file.GetStashForDate("2026-11-02")
	assert.NotNil(t, res)

	res = file.GetStashForDate("2026-11-03")
	assert.Nil(t, res)
}
