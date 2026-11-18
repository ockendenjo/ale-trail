package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"time"

	"github.com/ockendenjo/ale-trail/pkg/data"
	"github.com/ockendenjo/handler"
)

//go:embed overdue.txt
var overdueTemplate []byte

func (h *lambdaHandler) overdue(ctx *handler.Context, file *data.FileWithETag) (any, error) {
	revealDate := time.Now().AddDate(0, 0, 1)
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

	if maybeStash.Details != nil {
		logger.Info("Details already set - exiting")
		return nil, nil
	}

	tpl, err := template.New("overdue").Parse(string(overdueTemplate))
	if err != nil {
		return nil, err
	}

	data := overdueData{
		LocationLink: fmt.Sprintf("https://%s/%s/location?key=%s", h.domain, revealDateStr, maybeStash.StashKey),
	}

	buf := new(bytes.Buffer)
	if err = tpl.Execute(buf, data); err != nil {
		return nil, err
	}

	err = h.sendEmail(ctx, maybeStash.Stasher.Email, "Ale Trail - Add stash location now", buf.String())
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("Stash for %s does not have a location set", revealDateStr)
	err = h.sendTopicMessage(ctx, msg)
	return nil, err
}

type overdueData struct {
	LocationLink string
}
