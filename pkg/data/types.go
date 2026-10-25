package data

type File struct {
	Schema   *string `json:"$schema,omitempty"`
	AdminKey string  `json:"adminKey"`
	Stashes  []Stash `json:"stashes"`
}

type Stash struct {
	Date        string   `json:"date"`
	Details     *Details `json:"details,omitempty"`
	Image       string   `json:"image"`
	ReleaseTime string   `json:"releaseTime"`
	ReleaseType string   `json:"releaseType"`
	Stasher     *string  `json:"stasher,omitempty"`
	StashKey    string   `json:"stashKey"`
}

type Details struct {
	Contents    string  `json:"contents"`
	Description string  `json:"description"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	W3W         string  `json:"w3w"`
}
