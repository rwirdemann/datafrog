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

// The Verifier verifies the expectations in expectationSource. It monitors the
// databsaseLog for these expectations and requires them to be in same order as
// given in expectationSource. Verified expectations are written to
// verificationSink.
type Verifier struct {
	config            config.Config
	databsaseLog      ports.Log
	expectationSource ports.ExpectationSource
	timer             ports.Timer
	running           bool
}

// NewVerifier creates a new Verifier.
func NewVerifier(c config.Config, log ports.Log, source ports.ExpectationSource, t ports.Timer) *Verifier {
	return &Verifier{
		config:            c,
		databsaseLog:      log,
		expectationSource: source,
		timer:             t,
		running:           false,
	}
}

// Start runs the verification loop. It stops, when the expectations got out of
// order or when all expectations where met.
func (v *Verifier) Start() error {
	v.running = true
	v.timer.Start()
	log.Printf("Verification started at %v. Press Enter to stop verification...", v.timer.GetStart())
	expectations := v.expectationSource.GetAll()
	for i := range expectations {
		expectations[i].Fulfilled = false
	}
	for {
		if !v.running {
			v.expectationSource.WriteAll()
			break
		}
		line, err := v.databsaseLog.NextLine()
		if err != nil {
			log.Fatal(err)
		}

		// Hack to enable test adapter to stop the recording
		if line == "STOP" {
			v.expectationSource.WriteAll()
			break
		}

		ts, err := v.databsaseLog.Timestamp(line)
		if err != nil {
			continue
		}
		if v.timer.MatchesRecordingPeriod(ts) {
			matches, pattern := matcher.MatchesPattern(v.config, line)
			if matches {
				expectations := v.expectationSource.GetAll()
				for i, e := range expectations {
					if e.Fulfilled || e.Pattern != pattern {
						continue
					}

					if e.Verified == 0 {
						normalized := matcher.Normalize(line, v.config.Patterns)
						if diff, err := e.BuildDiff(normalized); err == nil {
							log.Printf("reference expectation found: %s\n", normalized)
							expectations[i].IgnoreDiffs = diff
							expectations[i].Fulfilled = true
							expectations[i].Verified = 1
							break
						}
					}

					if e.Verified > 0 {
						normalized := matcher.Normalize(line, v.config.Patterns)
						if e.Equal(normalized) {
							log.Printf("expectation verfied by: %s\n", normalized)
							expectations[i].Fulfilled = true
							expectations[i].Verified = e.Verified + 1
							break
						}
					}
				}
			}
		}
	}
	return nil
}

// Stop stops the verification.
func (v *Verifier) Stop() {
	v.running = false
	expectations := v.expectationSource.GetAll
	fulfilled := 0
	for _, e := range expectations() {
		if e.Fulfilled {
			fulfilled = fulfilled + 1
		}
	}
	log.Printf("Fulfilled %d of %d expectations\n", fulfilled, len(expectations()))
	for _, e := range expectations() {
		if !e.Fulfilled {
			log.Printf("Unfulfilled: '%s...'. Verification quote: %d", e.Tokens[0], e.Verified)
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
