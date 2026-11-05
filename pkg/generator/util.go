package generator

import (
	"fmt"
	"log"
	"math/rand/v2"
	"time"
	"uuid"

	"github.com/ockendenjo/ale-trail/pkg/data"
)

func NewFile(tm time.Time) data.File {
	file := data.File{
		Schema:   new("https://raw.githubusercontent.com/ockendenjo/ale-trail/refs/heads/main/schema/data.schema.json"),
		AdminKey: uuid.New().String(),
		Stashes:  getStashes(tm),
	}
	return file
}

func getStashes(tm time.Time) []*data.Stash {
	a := make([]*data.Stash, 0, 24)

	releaseTimes := makeReleaseTimes()

	for i := range 24 {
		releaseType := "RANDOM"
		if i%3 == 0 {
			releaseType = "TIMED"
		}

		stash := &data.Stash{
			Date:        tm.Format("2006-01-02"),
			ReleaseType: releaseType,
			ReleaseTime: releaseTimes[i],
			Image:       "",
			StashKey:    uuid.New().String(),
		}
		a = append(a, stash)
		tm = tm.AddDate(0, 0, 1)
	}

	return a
}

func makeReleaseTimes() []string {
	hourMap := map[int]int{}

	a := make([]string, 0, 24)

	for range 24 {
		var h, m int
		found := false

		for range 10 {
			h, m = getRandomReleaseTime()
			if hourMap[h] < 2 {
				hourMap[h]++
				found = true
				break
			}
		}

		if !found {
			log.Fatalf("infinite loop")
		}

		a = append(a, fmt.Sprintf("%02d:%02d", h, m))
	}

	return a
}

func getRandomReleaseTime() (int, int) {
	h := rand.IntN(18) + 3 // #nosec G404 - secure rand not needed
	m := rand.IntN(60)     // #nosec G404 - secure rand not needed
	return h, m
}
