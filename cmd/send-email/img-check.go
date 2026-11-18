package main

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/ockendenjo/ale-trail/pkg/data"
	"github.com/ockendenjo/handler"
)

func (h *lambdaHandler) imageCheck(ctx *handler.Context, file *data.FileWithETag) (any, error) {
	revealDate := time.Now().AddDate(0, 0, 1)
	revealDateStr := revealDate.Format(time.DateOnly)
	logger := ctx.GetLogger().AddParam("date", revealDateStr)

	maybeStash := file.GetStashForDate(revealDateStr)
	if maybeStash == nil {
		logger.Info("Stash not found")
		return nil, nil
	}

	if maybeStash.Stasher == nil {
		logger.Info("Stasher not set - exiting")
		return nil, nil
	}

	if maybeStash.Image == "" {
		logger.Info("Stash image not set")
		msg := fmt.Sprintf("Stash for %s does not have an image set", revealDateStr)
		err := h.sendTopicMessage(ctx, msg)
		return nil, err
	}

	if maybeStash.Image != "" {
		logger.Info("Stash image set")
	}

	ok, err := h.webClient.FileExists(ctx, maybeStash.Image)
	if err != nil {
		return nil, err
	}
	if !ok {
		logger.Info("Image does not exist in web bucket")
		err := h.sendTopicMessage(ctx, fmt.Sprintf("Image for %s does not exist in web bucket", revealDateStr))
		return nil, err
	}

	logger.Info("Image exists in web bucket")
	return nil, nil
}
