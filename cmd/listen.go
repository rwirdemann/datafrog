package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/rwirdemann/databasedragon/adapter"
	"github.com/rwirdemann/databasedragon/config"
	"github.com/rwirdemann/databasedragon/matcher"
	"github.com/rwirdemann/databasedragon/ports"
	"github.com/rwirdemann/databasedragon/ticker"
	"github.com/spf13/cobra"
)

func init() {
	listenCmd.Flags().String("expectations", "", "Filename with expectations")
	listenCmd.MarkFlagRequired("expectations")
	rootCmd.AddCommand(listenCmd)
}

type Listener struct {
	config             config.Config
	running            bool
	matcher            matcher.TokenMatcher
	databaseLog        ports.Log
	expectationSource  ports.ExpectationSource
	verificationSource ports.ExpectationSource
}

func NewListener(c config.Config, databseLog ports.Log, expectationSource ports.ExpectationSource,
	verificationSource ports.ExpectationSource) *Listener {
	return &Listener{
		config:             c,
		databaseLog:        databseLog,
		expectationSource:  expectationSource,
		verificationSource: expectationSource,
		running:            false}
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
	expectations := l.expectationSource.GetAll()
	verifications := l.expectationSource.GetAll()
	l.matcher = matcher.NewTokenMatcher(l.config, expectations, verifications)

	for {
		actual, err := l.databaseLog.NextLine()
		if err != nil {
			log.Fatal(err)
		}
		ts, err := l.databaseLog.Timestamp(actual)
		if err != nil {
			continue
		}
		if !t.MatchesRecordingPeriod(ts) {
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

		expectationSource := adapter.NewFileExpectationSource(expectations)
		verificationSource := adapter.NewFileExpectationSource(fmt.Sprintf("%s.verify", expectations))
		databaseLog := adapter.NewMYSQLLog(c.Filename)
		defer databaseLog.Close()

		listener = NewListener(c, databaseLog, expectationSource, verificationSource)
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
