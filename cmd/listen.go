package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/rwirdemann/databasedragon/adapter"
	"github.com/rwirdemann/databasedragon/config"
	"github.com/rwirdemann/databasedragon/matcher"
	"github.com/rwirdemann/databasedragon/ports"
	"github.com/spf13/cobra"
)

func init() {
	listenCmd.Flags().String("expectations", "", "Filename with expectations")
	listenCmd.MarkFlagRequired("expectations")
	rootCmd.AddCommand(listenCmd)
}

type Listener struct {
	config             config.Config
	timer              ports.Timer
	running            bool
	matcher            matcher.TokenMatcher
	databaseLog        ports.Log
	expectationSource  ports.ExpectationSource
	verificationSource ports.ExpectationSource
}

func NewListener(c config.Config, timer ports.Timer, databseLog ports.Log, expectationSource ports.ExpectationSource,
	verificationSource ports.ExpectationSource) *Listener {
	return &Listener{
		config:             c,
		timer:              timer,
		databaseLog:        databseLog,
		expectationSource:  expectationSource,
		verificationSource: verificationSource,
		running:            false}
}

// Start starts listening by checking each new logfile entry against the expectations from the
// expecations file. Matching expectations are removed. The listening counts as validated if all
// expectations were met and removed. The caller should stop the listening by calling
// Listener.Stop().
func (l *Listener) Start() {
	l.running = true
	l.timer.Start()
	log.Printf("Listening started at %v. Press Enter to stop listening...\n", l.timer.GetStart())
	expectations := l.expectationSource.GetAll()
	verifications := l.verificationSource.GetAll()
	l.matcher = matcher.NewTokenMatcher(l.config, expectations, verifications)

	for {
		actual, err := l.databaseLog.NextLine()
		if err != nil {
			log.Fatal(err)
		}
		if actual == "STOP" {
			break
		}

		ts, err := l.databaseLog.Timestamp(actual)
		if err != nil {
			continue
		}
		if !l.timer.MatchesRecordingPeriod(ts) {
			continue
		}

		if matchIndex := l.matcher.Matches(actual); matchIndex > -1 {
			l.matcher.RemoveExpectation(matchIndex)
		}
	}
}

// Stop stops the listening and validation loop.
func (l *Listener) Stop() {
	log.Println("Listening stoped")
	l.matcher.PrintResults()
}

func (l *Listener) GetResults() []matcher.Expectation {
	return l.matcher.GetResults()
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

		t := &adapter.UTCTimer{}
		expectationSource := adapter.NewFileExpectationSource(expectations)
		verificationSource := adapter.NewFileExpectationSource(fmt.Sprintf("%s.verify", expectations))
		databaseLog := adapter.NewMYSQLLog(c.Filename)
		defer databaseLog.Close()

		listener = NewListener(c, t, databaseLog, expectationSource, verificationSource)
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
