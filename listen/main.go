package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rwirdemann/databasedragon/adapter"
	"github.com/rwirdemann/databasedragon/config"
	"github.com/rwirdemann/databasedragon/matcher"
	"github.com/rwirdemann/databasedragon/ticker"
	"github.com/rwirdemann/databasedragon/validation"
)

var validator validation.Validator

func main() {
	c := config.NewConfig("config.json")
	println(c.Filename)
	fmt.Print("Press Enter to start listening...")
	_, _ = fmt.Scanln()

	t := ticker.Ticker{}
	t.Start()
	fmt.Printf("Listening started at  %v. Press Enter to stop recording...\n", t.GetStart())

	validationFile, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	validator = validation.NewUnorderedRemovalValidator(strings.Split(string(validationFile), "\n"))

	logPort := adapter.NewMYSQLLog(c.Filename)
	defer logPort.Close()

	go checkExit()

	m := matcher.NewDynamicDataMatcher(c)
	for {
		line, err := logPort.NextLine()
		if err != nil {
			log.Fatal(err)
		}
		ts, err := logPort.Timestamp(line)
		if err != nil {
			continue
		}
		if t.MatchesRecordingPeriod(ts) {
			matches := false
			matchingExpectations := ""
			if m.MatchesAny(line) {
				for _, e := range validator.GetExpectations() {
					if m.MatchesExactly(line, e) {
						fmt.Printf("Expectation met: %s", line)
						matches = true
						matchingExpectations = e
						break
					}
				}
			}
			if matches {
				validator.Remove(matchingExpectations)
			}
		}
	}
}

func checkExit() {
	var b = make([]byte, 1)
	l, _ := os.Stdin.Read(b)
	if l > 0 {
		validator.PrintResults()
		os.Exit(0)
	}
}
