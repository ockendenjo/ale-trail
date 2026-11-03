package api

import (
	"time"

	"github.com/ockendenjo/ale-trail/pkg/data"
)

type Stash struct {
	Date        string `json:"date"`
	Locked      bool   `json:"locked"`
	ReleaseType string `json:"releaseType"`

	//Maybe when locked
	HostNeeded  *bool   `json:"hostNeeded,omitempty"`
	ReleaseTime *string `json:"releaseTime,omitempty"`

	//When unlocked
	Image       *string  `json:"image,omitempty"`
	Contents    *string  `json:"contents,omitempty"`
	Description *string  `json:"description,omitempty"`
	Lat         *float64 `json:"lat,omitempty"`
	Lon         *float64 `json:"lon,omitempty"`
	W3W         *string  `json:"w3w,omitempty"`
}

type Stashes struct {
	Stashes []*Stash `json:"stashes"`
}

func GetStash(data *data.Stash, fullMode bool) *Stash {

	locked := data.IsLocked(time.Now())

	stash := &Stash{
		Date:        data.Date,
		Locked:      locked,
		ReleaseType: data.ReleaseType,
	}

	if locked {
		stash.HostNeeded = new(data.Stasher == nil)
		if data.ReleaseType == "TIMED" && data.IsToday() {
			stash.ReleaseTime = &data.ReleaseTime
		}
	}

	if !fullMode {
		return stash
	}

	if !locked {
		details := data.Details
		stash.W3W = &details.W3W
		stash.Lon = &details.Lon
		stash.Lat = &details.Lat
		stash.Description = &details.Description
		stash.Contents = &details.Contents
		stash.Image = &data.Image
	}

	return stash
}
