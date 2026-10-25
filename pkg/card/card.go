package card

type Card struct {
	Date        string  `json:"date"`
	Locked      bool    `json:"locked"`
	ReleaseType *string `json:"releaseType"`
	ReleaseTime *string `json:"releaseTime"`
	Image       *string `json:"image"`
	HostNeeded  *bool   `json:"hostNeeded"`
}
