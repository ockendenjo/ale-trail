package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ockendenjo/ale-trail/pkg/generator"
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

	file := generator.NewFile(tm)

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(file)
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}
