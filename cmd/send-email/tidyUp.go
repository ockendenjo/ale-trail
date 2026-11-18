package main

import (
	_ "embed"
	"time"

	"github.com/ockendenjo/ale-trail/pkg/data"
	"github.com/ockendenjo/handler"
)

//go:embed tidyUp.txt
var tidyUpTemplate []byte

func (h *lambdaHandler) tidyUp(ctx *handler.Context, file *data.FileWithETag) (any, error) {
	revealDate := time.Now().AddDate(0, 0, -1)
	revealDateStr := revealDate.Format(time.DateOnly)
	logger := ctx.GetLogger().AddParam("date", revealDateStr)

	maybeStash := file.GetStashForDate(revealDateStr)
	if maybeStash == nil {
		logger.Info("Stash not found")
		return nil, nil
	}

	if maybeStash.Stasher == nil {
		logger.Info("Stasher not set - ignoring")
		return nil, nil
	}

	if maybeStash.Details == nil {
		logger.Info("Details not set - ignoring")
		return nil, nil
	}

	err := h.sendEmail(ctx, maybeStash.Stasher.Email, "Ale Trail - Ensure stash removed", string(tidyUpTemplate))
	return nil, err
}
