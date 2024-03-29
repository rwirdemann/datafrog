package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/rwirdemann/databasedragon/adapter"
	"github.com/rwirdemann/databasedragon/config"
	"github.com/rwirdemann/databasedragon/matcher"
	"github.com/rwirdemann/databasedragon/ticker"
	"github.com/rwirdemann/databasedragon/validation"
)

type Listener struct {
	config              config.Config
	expectationFilename string
	running             bool
	validator           validation.Validator
}

func NewListener(c config.Config, expectationFilename string) *Listener {
	return &Listener{config: c, expectationFilename: expectationFilename, running: false}
}

func (l *Listener) Start() {
	l.running = true
	t := ticker.Ticker{}
	t.Start()
	log.Printf("Listening started at %v. Press Enter to stop listening...\n", t.GetStart())
	expecations, err := os.ReadFile(l.expectationFilename)
	if err != nil {
		log.Fatal(err)
	}
	l.validator = validation.NewUnorderedRemovalValidator(strings.Split(string(expecations), "\n"))

	logPort := adapter.NewMYSQLLog(l.config.Filename)
	defer logPort.Close()

	m := matcher.NewDynamicDataMatcher(l.config)
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
			matchingExpectation := ""
			if m.MatchesAny(line) {
				for _, e := range l.validator.GetExpectations() {
					if m.MatchesExactly(line, e) {
						log.Printf("Expectation met...: %s", line[:60])
						matches = true
						matchingExpectation = e
						break
					}
				}
			}
			if matches {
				l.validator.Remove(matchingExpectation)
			}
		}
	}
}

func (l *Listener) Stop() {
	log.Println("Listening stoped")
	l.validator.PrintResults()
}
