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

//go:embed dayMinus2.txt
var dayMinus2Template []byte

func (h *lambdaHandler) dayMinus2(ctx *handler.Context, file *data.FileWithETag) (any, error) {
	date := time.Now().AddDate(0, 0, 2)
	dateStr := date.Format(time.DateOnly)
	maybeStash := file.GetStashForDate(dateStr)
	if maybeStash == nil {
		ctx.GetLogger().Info("Stash not found")
		return nil, nil
	}

	if maybeStash.Stasher == nil {
		err := h.sendTopicMessage(ctx, fmt.Sprintf("Stash for %s does not have a stasher - not sending email", dateStr))
		return nil, err
	}

	tpl, err := template.New("dayMinus2").Parse(string(dayMinus2Template))
	if err != nil {
		return nil, err
	}

	data := dayMinus2Data{
		StashDOM:     date.Format("02"),
		StashDOW:     date.Format("Mon"),
		StashMonth:   date.Format("Jan"),
		LocationLink: fmt.Sprintf("https://%s/%s/location?key=%s", h.domain, dateStr, maybeStash.StashKey),
	}

	buf := new(bytes.Buffer)
	if err = tpl.Execute(buf, data); err != nil {
		return nil, err
	}

	err = h.sendEmail(ctx, maybeStash.Stasher.Email, "Ale Trail - Hide your stash tomorrow", buf.String())
	return nil, err
}

type dayMinus2Data struct {
	StashDOM         string
	StashDOW         string
	StashMonth       string
	RequiredContents string
	LocationLink     string
}
