package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"time"
	"uuid"

	"github.com/ockendenjo/ale-trail/pkg/data"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gen-data-file startDate")
		os.Exit(1)
	}
	startDate := os.Args[1]
	tm, err := time.Parse(time.DateOnly, startDate)
	if err != nil {
		fmt.Println("Invalid start date")
		os.Exit(1)
	}

	file := data.File{
		Schema:   new("https://raw.githubusercontent.com/ockendenjo/ale-trail/refs/heads/main/schema/data.schema.json"),
		AdminKey: uuid.New().String(),
		Stashes:  getStashes(tm),
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(file)
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}

func getStashes(tm time.Time) []data.Stash {
	a := make([]data.Stash, 0, 24)

	releaseTimes := makeReleaseTimes()

	for i := range 24 {
		releaseType := "RANDOM"
		if i%3 == 0 {
			releaseType = "TIMED"
		}

		stash := data.Stash{
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
	h := rand.IntN(18) + 3
	m := rand.IntN(60)
	return h, m
}
