package data

import "time"

type File struct {
	Schema   *string  `json:"$schema,omitempty"`
	AdminKey string   `json:"adminKey"`
	Stashes  []*Stash `json:"stashes"`
}

type Stash struct {
	Date        string   `json:"date"`
	Details     *Details `json:"details,omitempty"`
	Image       string   `json:"image"`
	ReleaseTime string   `json:"releaseTime"`
	ReleaseType string   `json:"releaseType"`
	Stasher     *Stasher `json:"stasher,omitempty"`
	StashKey    string   `json:"stashKey"`
}

func (s *Stash) IsToday() bool {
	return s.Date == time.Now().Format(time.DateOnly)
}

type Stasher struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Details struct {
	Contents    string  `json:"contents"`
	Description string  `json:"description"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	W3W         string  `json:"w3w"`
}

type FileWithETag struct {
	File *File   `json:"file"`
	ETag *string `json:"etag"`
}
