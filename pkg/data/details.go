package data

import (
	"errors"
	"strings"
)

type Details struct {
	Contents    string  `json:"contents"`
	Description string  `json:"description"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	W3W         string  `json:"w3w"`
}

func (d *Details) Validate() error {
	if d.Lat > 90 || d.Lat < -90 {
		return errors.New("invalid latitude value")
	}
	if d.Lon > 180 || d.Lon < -180 {
		return errors.New("invalid longitude value")
	}
	if d.W3W == "" {
		return errors.New("missing w3w value")
	}
	if after, ok := strings.CutPrefix(d.W3W, "///"); ok {
		d.W3W = after
	}
	if d.Description == "" {
		return errors.New("missing description value")
	}
	if len(d.Description) > 200 {
		return errors.New("description too long")
	}
	if len(d.Contents) > 200 {
		return errors.New("content too long")
	}
	return nil
}
