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

//go:embed dayMinus7.txt
var dayMinus7Template []byte

func (h *lambdaHandler) dayMinus7(ctx *handler.Context, file *data.FileWithETag) (any, error) {
	revealDate := time.Now().AddDate(0, 0, 7)
	revealDateStr := revealDate.Format(time.DateOnly)

	maybeStash := file.GetStashForDate(revealDateStr)
	if maybeStash == nil {
		ctx.GetLogger().Info("Stash not found")
		return nil, nil
	}

	if maybeStash.Stasher == nil {
		ctx.GetLogger().Info("Stasher not set - exiting")
		return nil, nil
	}

	text, err := h.getDayMinus7Text(revealDate, maybeStash)
	if err != nil {
		return nil, err
	}

	err = h.sendEmail(ctx, maybeStash.Stasher.Email, "Ale Trail - 7 day reminder", text)
	return nil, err
}

func (h *lambdaHandler) getDayMinus7Text(revealDate time.Time, stash *data.Stash) (string, error) {
	tpl, err := template.New("dayMinus7").Parse(string(dayMinus7Template))
	if err != nil {
		return "", err
	}

	revealDateStr := revealDate.Format(time.DateOnly)

	stashDate := revealDate.AddDate(0, 0, -1)

	data := dayMinus7Data{
		RevealDOM:   revealDate.Format("02"),
		RevealDOW:   revealDate.Format("Mon"),
		RevealMonth: revealDate.Format("Jan"),

		StashDOM:    stashDate.Format("02"),
		StashDOW:    stashDate.Format("Mon"),
		StashMonth:  stashDate.Format("Jan"),
		StasherLink: fmt.Sprintf("https://%s/%s/stasher?key=%s", h.domain, revealDateStr, stash.StashKey),
	}

	buf := new(bytes.Buffer)
	if err = tpl.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

type dayMinus7Data struct {
	RevealDOM   string
	RevealDOW   string
	RevealMonth string

	StashDOM   string
	StashDOW   string
	StashMonth string

	StasherLink string
}
