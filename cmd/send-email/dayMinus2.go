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
	revealDate := time.Now().AddDate(0, 0, 2)
	revealDateStr := revealDate.Format(time.DateOnly)
	logger := ctx.GetLogger().AddParam("date", revealDateStr)

	maybeStash := file.GetStashForDate(revealDateStr)
	if maybeStash == nil {
		logger.Info("Stash not found")
		return nil, nil
	}

	if maybeStash.Stasher == nil {
		err := h.sendTopicMessage(ctx, fmt.Sprintf("Stash for %s does not have a stasher - not sending email", revealDateStr))
		return nil, err
	}

	emailBody, err := h.getDayMinus2Text(revealDate, maybeStash)
	if err != nil {
		return nil, err
	}

	err = h.sendEmail(ctx, maybeStash.Stasher.Email, "Ale Trail - Hide your stash tomorrow", emailBody)
	return nil, err
}

func (h *lambdaHandler) getDayMinus2Text(revealDate time.Time, stash *data.Stash) (string, error) {
	tpl, err := template.New("dayMinus2").Parse(string(dayMinus2Template))
	if err != nil {
		return "", err
	}

	revealDateStr := revealDate.Format(time.DateOnly)
	hideDate := revealDate.AddDate(0, 0, -1)

	data := dayMinus2Data{
		StashDOM:     hideDate.Format("02"),
		StashDOW:     hideDate.Format("Mon"),
		StashMonth:   hideDate.Format("Jan"),
		LocationLink: fmt.Sprintf("https://%s/%s/location?key=%s", h.domain, revealDateStr, stash.StashKey),
	}

	buf := new(bytes.Buffer)
	if err = tpl.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

type dayMinus2Data struct {
	StashDOM         string
	StashDOW         string
	StashMonth       string
	RequiredContents string
	LocationLink     string
}
