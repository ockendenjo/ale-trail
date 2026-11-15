package main

import (
	"testing"
	"time"

	"github.com/ockendenjo/ale-trail/pkg/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getDayMinus1Text(t *testing.T) {
	h := &lambdaHandler{
		domain: "example.com",
	}

	revealDate := time.Date(2026, 11, 16, 0, 0, 0, 0, time.UTC)
	stash := &data.Stash{
		StashKey: "secret",
	}

	text, err := h.getDayMinus1Text(revealDate, stash)
	require.NoError(t, err)

	assert.Contains(t, text, "Today (Sun 15 Nov)")
	assert.Contains(t, text, "https://example.com/2026-11-16/location?key=secret")
}
