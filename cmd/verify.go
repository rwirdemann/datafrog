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
	verifyCmd.Flags().String("expectations", "", "Filename to save verify")
	verifyCmd.Flags().Bool("prompt", false, "Wait for key stroke before verification starts")
	verifyCmd.MarkFlagRequired("expectations")
	rootCmd.AddCommand(verifyCmd)
}

// The Verifier verifies the initialExpectations in expectationSource. It monitors the
// databaseLog for these initialExpectations and requires them to be in same order as
// given in expectationSource. Verified initialExpectations are written to
// verificationSink.
type Verifier struct {
	config            config.Config
	databaseLog       ports.Log
	expectationSource ports.ExpectationSource
	timer             ports.Timer
	running           bool
}

// NewVerifier creates a new Verifier.
func NewVerifier(c config.Config, log ports.Log, source ports.ExpectationSource, t ports.Timer) *Verifier {
	return &Verifier{
		config:            c,
		databaseLog:       log,
		expectationSource: source,
		timer:             t,
		running:           false,
	}
}

// Start runs the verification loop. It stops, when the initialExpectations got out of
// order or when all initialExpectations where met.
func (v *Verifier) Start() error {
	v.running = true
	v.timer.Start()
	log.Printf("Verification started at %v. Press Enter to stop and save verification...", v.timer.GetStart())
	expectations := v.expectationSource.GetAll()
	for i := range expectations {
		expectations[i].Fulfilled = false
	}
	for {
		if !v.running {
			v.expectationSource.WriteAll()
			break
		}
		line, err := v.databaseLog.NextLine()
		if err != nil {
			log.Fatal(err)
		}

		// Hack to enable test adapter to stop the recording
		if line == "STOP" {
			v.expectationSource.WriteAll()
			break
		}

		ts, err := v.databaseLog.Timestamp(line)
		if err != nil {
			continue
		}
		if v.timer.MatchesRecordingPeriod(ts) {
			matches, pattern := matcher.MatchesPattern(v.config, line)
			if !matches {
				continue
			}

			expectations := v.expectationSource.GetAll()
			for i, e := range expectations {
				if e.Fulfilled || e.Pattern != pattern {
					continue
				}

				tokens := matcher.Tokenize(matcher.Normalize(line, v.config.Patterns))

				// handle already verified expectations
				if e.Verified > 0 && e.Equal(tokens) {
					log.Printf("expectation verified by: %s\n", matcher.Expectation{Tokens: tokens}.Shorten(6))
					expectations[i].Fulfilled = true
					expectations[i].Verified = e.Verified + 1
					break
				}

				// handle not yet verified expectations (verified == 0)
				if diff, err := e.Diff(tokens); err == nil {
					log.Printf("reference expectation found: %s\n", matcher.Expectation{Tokens: tokens}.Shorten(6))
					expectations[i].IgnoreDiffs = diff
					expectations[i].Fulfilled = true
					expectations[i].Verified = 1
					break
				}
			}
		}
		if allFulfilled(expectations) {
			v.Stop()
		}
	}
	return nil
}

func allFulfilled(expectations []matcher.Expectation) bool {
	for _, e := range expectations {
		if !e.Fulfilled {
			return false
		}
	}
	return true
}

// Stop stops the verification.
func (v *Verifier) Stop() {
	v.running = false
	expectations := v.expectationSource.GetAll
	fulfilled := 0
	verifiedSum := 0
	for _, e := range expectations() {
		verifiedSum += e.Verified
		if e.Fulfilled {
			fulfilled = fulfilled + 1
		}
	}
	log.Printf("Fulfilled %d of %d expectations\n", fulfilled, len(expectations()))
	log.Printf("Verification mean: %f\n", float32(verifiedSum)/float32(len(expectations())))
	for _, e := range expectations() {
		if !e.Fulfilled {
			log.Printf("Unfulfilled: '%s'. Verification quote: %d", e.Shorten(6), e.Verified)
		}
	}
}

var verifier *Verifier
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Starts verification",
	RunE: func(cmd *cobra.Command, args []string) error {
		expectationsFilename, _ := cmd.Flags().GetString("expectations")
		c := config.NewConfig("config.json")
		prompt, _ := cmd.Flags().GetBool("prompt")
		if prompt {
			log.Printf("Verifying '%s'. Hit enter when you are ready!", expectationsFilename)
			_, _ = fmt.Scanln()
		} else {
			log.Printf("Verifying '%s'.", expectationsFilename)
		}
		go checkVerifyExit()

		expectationSource := adapter.NewFileExpectationSource(expectationsFilename)
		databaseLog := adapter.NewMYSQLLog(c.Filename)
		defer databaseLog.Close()

		t := &adapter.UTCTimer{}

		verifier = NewVerifier(c, databaseLog, expectationSource, t)
		return verifier.Start()
	},
}

// Checks if enter was hit to stop verification.
func checkVerifyExit() {
	var b = make([]byte, 1)
	l, _ := os.Stdin.Read(b)
	if l > 0 {
		verifier.Stop()
	}
}
