package card

type Card struct {
	Date        string  `json:"date"`
	Locked      bool    `json:"locked"`
	ReleaseType *string `json:"releaseType,omitempty"`
	ReleaseTime *string `json:"releaseTime,omitempty"`
	Image       *string `json:"image,omitempty"`
	HostNeeded  *bool   `json:"hostNeeded,omitempty"`
}
