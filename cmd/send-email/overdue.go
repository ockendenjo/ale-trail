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
	date := time.Now().AddDate(0, 0, 1)
	dateStr := date.Format(time.DateOnly)
	maybeStash := file.GetStashForDate(dateStr)
	if maybeStash == nil {
		ctx.GetLogger().Info("Stash not found")
		return nil, nil
	}

	if maybeStash.Stasher == nil {
		ctx.GetLogger().Info("Stasher not set - ignoring")
		return nil, nil
	}

	if maybeStash.Details != nil {
		ctx.GetLogger().Info("Details already set - exiting")
		return nil, nil
	}

	tpl, err := template.New("overdue").Parse(string(overdueTemplate))
	if err != nil {
		return nil, err
	}

	data := overdueData{
		LocationLink: fmt.Sprintf("https://%s/%s/location?key=%s", h.domain, dateStr, maybeStash.StashKey),
	}

	buf := new(bytes.Buffer)
	if err = tpl.Execute(buf, data); err != nil {
		return nil, err
	}

	err = h.sendEmail(ctx, maybeStash.Stasher.Email, "Ale Trail - Add stash location now", buf.String())
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("Stash for %s does not have a location set", dateStr)
	err = h.sendTopicMessage(ctx, msg)
	return nil, err
}

type overdueData struct {
	LocationLink string
}
