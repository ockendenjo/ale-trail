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

//go:embed dayMinus1.txt
var dayMinus1Template []byte

func (h *lambdaHandler) dayMinus1(ctx *handler.Context, file *data.FileWithETag) (any, error) {
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

	tpl, err := template.New("dayMinus2").Parse(string(dayMinus1Template))
	if err != nil {
		return nil, err
	}

	data := dayMinus1Data{
		StashDOM:     date.Format("02"),
		StashDOW:     date.Format("Mon"),
		StashMonth:   date.Format("Jan"),
		LocationLink: fmt.Sprintf("https://%s/%s/location?key=%s", h.domain, dateStr, maybeStash.StashKey),
	}

	buf := new(bytes.Buffer)
	if err = tpl.Execute(buf, data); err != nil {
		return nil, err
	}

	err = h.sendEmail(ctx, maybeStash.Stasher.Email, "Ale Trail - Hide your stash today", buf.String())
	return nil, err
}

type dayMinus1Data struct {
	StashDOM     string
	StashDOW     string
	StashMonth   string
	LocationLink string
}
