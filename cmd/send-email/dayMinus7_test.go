package main

import (
	_ "embed"
	"testing"
	"time"

	"github.com/ockendenjo/ale-trail/pkg/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed dayMinus7.output.txt
var expDayMinus7Text []byte

func Test_getDayMinus7Text(t *testing.T) {
	h := &lambdaHandler{
		domain: "example.com",
	}

	revealDate := time.Date(2026, 11, 25, 0, 0, 0, 0, time.UTC)
	stash := &data.Stash{
		StashKey: "secret",
	}

	text, err := h.getDayMinus7Text(revealDate, stash)
	require.NoError(t, err)

	assert.Equal(t, string(expDayMinus7Text), text)
}
