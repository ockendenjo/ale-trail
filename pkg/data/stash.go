package data

import "time"

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

func (s *Stash) IsLocked(instant time.Time) bool {
	return !s.IsUnlocked(instant)
}

func (s *Stash) IsUnlocked(instant time.Time) bool {
	if s.Details == nil {
		return false
	}
	td := s.Date + "T" + s.ReleaseTime + ":00Z"
	tm, err := time.Parse("2006-01-02T15:04:05Z", td)
	if err != nil {
		return false
	}
	return tm.Before(instant)
}

type Stasher struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
