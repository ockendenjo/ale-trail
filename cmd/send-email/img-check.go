package main

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/ockendenjo/ale-trail/pkg/data"
	"github.com/ockendenjo/handler"
)

func (h *lambdaHandler) imageCheck(ctx *handler.Context, file *data.FileWithETag) (any, error) {
	date := time.Now().AddDate(0, 0, 1)
	dateStr := date.Format(time.DateOnly)
	maybeStash := file.GetStashForDate(dateStr)
	if maybeStash == nil {
		ctx.GetLogger().Info("Stash not found")
		return nil, nil
	}

	if maybeStash.Stasher == nil {
		ctx.GetLogger().Info("Stasher not set - exiting")
		return nil, nil
	}

	if maybeStash.Image != "" {
		ctx.GetLogger().Info("Image set - exiting")
		return nil, nil
	}

	msg := fmt.Sprintf("Stash for %s does not have an image set", dateStr)
	err := h.sendTopicMessage(ctx, msg)
	return nil, err
}
