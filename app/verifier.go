package app

import (
	"fmt"
	"github.com/rwirdemann/databasedragon/app/domain"
	"github.com/rwirdemann/databasedragon/config"
	"github.com/rwirdemann/databasedragon/matcher"
	"github.com/rwirdemann/databasedragon/ports"
	"log"
	"time"
)

// The Verifier verifies the initialExpectations in expectationSource. It
// monitors the databaseLog for these initialExpectations and requires them to
// be in same order as given in expectationSource. Verified initialExpectations
// are written to verificationSink.
type Verifier struct {
	config            config.Config
	tokenizer         matcher.Tokenizer
	databaseLog       ports.Log
	expectationSource ports.ExpectationSource
	testcase          domain.Testcase
	timer             ports.Timer
	name              string
}

// NewVerifier creates a new Verifier.
func NewVerifier(c config.Config, tokenizer matcher.Tokenizer, log ports.Log,
	source ports.ExpectationSource, t ports.Timer, name string) *Verifier {
	return &Verifier{
		config:            c,
		tokenizer:         tokenizer,
		databaseLog:       log,
		expectationSource: source,
		testcase:          source.Get(),
		timer:             t,
		name:              name,
	}
}

// Start runs the verification loop.
func (verifier *Verifier) Start(done chan struct{}, stopped chan struct{}) {
	verifier.timer.Start()
	log.Printf("Verification started at %v. Press Enter to stop and save verification...", verifier.timer.GetStart())
	verifier.testcase.Verifications = verifier.testcase.Verifications + 1
	for i := range verifier.testcase.Expectations {
		verifier.testcase.Expectations[i].Fulfilled = false
	}

	// tell caller that verification has been finished
	defer close(stopped)

	// called when done channel is closed
	defer func() {
		_ = verifier.expectationSource.Write(verifier.testcase)
	}()

	for {
		select {
		default:
			if allFulfilled(verifier.testcase.Expectations) {
				log.Printf("All verifications fulfilled. Verification done")
				return
			}

			v, err := verifier.databaseLog.NextLine()
			if err != nil {
				log.Fatal(err)
			}

			ts, err := verifier.databaseLog.Timestamp(v)
			if err != nil {
				continue
			}
			if verifier.timer.MatchesRecordingPeriod(ts) {
				matches, pattern := matcher.MatchesPattern(verifier.config, v)
				if !matches {
					continue
				}

				for i, e := range verifier.testcase.Expectations {
					if e.Fulfilled || e.Pattern != pattern {
						continue // -> continue with next e
					}

					vTokens := verifier.tokenizer.Tokenize(v, verifier.config.Patterns)

					// Handle already verified expectations.
					if e.Verified > 0 && e.Equal(vTokens) {
						log.Printf("expectation verified by: %s\n", domain.Expectation{Tokens: vTokens}.Shorten(6))
						verifier.testcase.Expectations[i].Fulfilled = true
						verifier.testcase.Expectations[i].Verified = e.Verified + 1
						break // -> continue with next verifier
					}

					if len(e.Tokens) != len(vTokens) {
						break // -> continue with next verifier
					}

					// Not yet verified expectation e with same token lengths as
					// verifier found. This expectation e becomes our references
					// expectation.
					if diff, err := e.Diff(vTokens); err == nil {
						log.Printf("reference expectation found: %s\n", domain.Expectation{Tokens: vTokens}.Shorten(6))
						verifier.testcase.Expectations[i].IgnoreDiffs = diff
						verifier.testcase.Expectations[i].Fulfilled = true
						verifier.testcase.Expectations[i].Verified = 1
						break // -> continue with next verifier
					}
				}
			}
		case <-done:
			log.Printf("Channel close: Verification done")
			return
		}
	}
}

// allFulfilled checks all expectations, returns true if all fulfilled and false
// otherwise.
func allFulfilled(expectations []domain.Expectation) bool {
	for _, e := range expectations {
		if !e.Fulfilled {
			return false
		}
	}
	return true
}

// ReportResults reports the verification results.
func (verifier *Verifier) ReportResults() domain.Report {
	fulfilled := 0
	verifiedSum := 0
	for _, e := range verifier.testcase.Expectations {
		verifiedSum += e.Verified
		if e.Fulfilled {
			fulfilled = fulfilled + 1
		}
	}
	report := domain.Report{
		Testname:         verifier.name,
		LastExecution:    time.Now(),
		Expectations:     len(verifier.testcase.Expectations),
		Verifications:    verifier.testcase.Verifications,
		Fulfilled:        fulfilled,
		VerificationMean: verificationMean(float32(verifiedSum), float32(len(verifier.testcase.Expectations))),
	}
	for _, e := range verifier.testcase.Expectations {
		if !e.Fulfilled {
			report.Unfulfilled = append(report.Unfulfilled, fmt.Sprintf("%s. Verification quote: %d", e.Shorten(6), e.Verified))
		}
	}
	return report
}

func verificationMean(sum, expectationCount float32) float32 {
	if expectationCount > 0 {
		return sum / expectationCount
	}
	return 0
}
