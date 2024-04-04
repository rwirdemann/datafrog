package cmd

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
	"github.com/spf13/cobra"
)

func init() {
	listenCmd.Flags().String("expectations", "", "Filename with expectations")
	listenCmd.MarkFlagRequired("expectations")
	rootCmd.AddCommand(listenCmd)
}

type Listener struct {
	config              config.Config
	expectationFilename string
	running             bool
	validator           validation.Validator
}

func NewListener(c config.Config, expectationFilename string) *Listener {
	return &Listener{config: c, expectationFilename: expectationFilename, running: false}
}

// Start starts listening by checking each new logfile entry against the expectations from the
// expecations file. Matching expectations are removed. The listening counts as validated if all
// expectations were met and removed. The caller should stop the listening by calling
// Listener.Stop().
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
	initialExpectationCount := len(l.validator.GetExpectations())

	logPort := adapter.NewMYSQLLog(l.config.Filename)
	defer logPort.Close()

	m := matcher.NewLevenshteinMatcher(l.config)
	for {
		line, err := logPort.NextLine()
		if err != nil {
			log.Fatal(err)
		}
		ts, err := logPort.Timestamp(line)
		if err != nil {
			continue
		}
		if !t.MatchesRecordingPeriod(ts) {
			continue
		}

		var Green = "\033[32m"
		var Yellow = "\033[33m"
		var White = "\033[97m"
		log.SetFlags(0)
		matchesPattern, pattern := m.MatchesPattern(line)
		if matchesPattern {
			log.Printf(Green+"PATTERN: %v", pattern)
			log.Printf(Green+"IN: %s", line)
			for i, e := range l.validator.GetExpectations() {
				if pattern.MatchesAllConditions(e) {
					log.Printf(Yellow+"E%d: %s", i+1, e)
					if m.MatchesExactly(line, e) {
						l.validator.Remove(e)
						log.Printf(White+"Remaining Exceptions: %d / %d\n", len(l.validator.GetExpectations()), initialExpectationCount)
						break
					}
				}
			}
		}
	}
}

// Stop stops the listening and validation loop.
func (l *Listener) Stop() {
	log.Println("Listening stoped")
	l.validator.PrintResults()
}

var listener *Listener
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Starts listening and validation",
	RunE: func(cmd *cobra.Command, args []string) error {
		expectations, _ := cmd.Flags().GetString("expectations")
		c := config.NewConfig("config.json")
		log.Printf("Listening to '%s'. Hit enter when you are ready!", expectations)
		_, _ = fmt.Scanln()
		go checkStopListening()

		listener = NewListener(c, expectations)
		listener.Start()

		return nil
	},
}

// Checks if enter was hit to stop listening.
func checkStopListening() {
	var b = make([]byte, 1)
	l, _ := os.Stdin.Read(b)
	if l > 0 {
		listener.Stop()
		os.Exit(0)
	}
}
