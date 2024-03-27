package main

import (
	"fmt"
	"github.com/rwirdemann/databasedragon/adapter"
	"github.com/rwirdemann/databasedragon/config"
	"github.com/rwirdemann/databasedragon/matcher"
	"github.com/rwirdemann/databasedragon/ticker"
	"github.com/rwirdemann/databasedragon/validation"
	"log"
	"os"
	"strings"
	"time"
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

	m := matcher.NewPatternMatcher(c)
	for {
		line, err := logPort.NextLine()
		if err != nil {
			log.Fatal(err)
		}
		ts, err := logPort.Timestamp(line)
		if err != nil {
			continue
		}
		if t.MatchesRecordingPeriod(ts) && m.MatchesAny(line) {
			fmt.Printf("Expectation met: %s", line)
			pattern := m.MatchingPattern(line)
			validator.RemoveFirstMatchingExpectation(pattern)
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
func matchesRecordingPeriod(ts time.Time, startDate time.Time) bool {
	return ts.Equal(startDate) || ts.After(startDate)
}
