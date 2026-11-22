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

//go:embed dayMinus1.template.txt
var dayMinus1Template []byte

func (h *lambdaHandler) dayMinus1(ctx *handler.Context, file *data.FileWithETag) (any, error) {
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

	text, err := h.getDayMinus1Text(revealDate, maybeStash)
	if err != nil {
		return nil, err
	}

	err = h.sendEmail(ctx, maybeStash.Stasher.Email, "Ale Trail - Hide your stash today", text)
	return nil, err
}

func (h *lambdaHandler) getDayMinus1Text(revealDate time.Time, stash *data.Stash) (string, error) {
	tpl, err := template.New("dayMinus1").Parse(string(dayMinus1Template))
	if err != nil {
		return "", err
	}

	revealDateStr := revealDate.Format(time.DateOnly)

	hideDate := revealDate.AddDate(0, 0, -1)

	data := dayMinus1Data{
		StashDOM:     hideDate.Format("02"),
		StashDOW:     hideDate.Format("Mon"),
		StashMonth:   hideDate.Format("Jan"),
		LocationLink: fmt.Sprintf("https://%s/%s/location?key=%s", h.domain, revealDateStr, stash.StashKey),
		PreviousLink: fmt.Sprintf("https://%s/previous-locations", h.domain),
	}

	buf := new(bytes.Buffer)
	if err = tpl.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

type dayMinus1Data struct {
	StashDOM     string
	StashDOW     string
	StashMonth   string
	LocationLink string
	PreviousLink string
}
