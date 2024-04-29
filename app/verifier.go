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
func (v *Verifier) Start(done chan struct{}, stopped chan struct{}) {
	v.timer.Start()
	log.Printf("Verification started at %v. Press Enter to stop and save verification...", v.timer.GetStart())
	v.testcase.Verifications = v.testcase.Verifications + 1
	for i := range v.testcase.Expectations {
		v.testcase.Expectations[i].Fulfilled = false
	}

	// tell caller that verification has been finished
	defer close(stopped)

	// called when done channel is closed
	defer func() {
		_ = v.expectationSource.Write(v.testcase)
	}()

	for {
		select {
		default:
			if allFulfilled(v.testcase.Expectations) {
				log.Printf("All verifications fulfilled. Verification done")
				return
			}

			line, err := v.databaseLog.NextLine()
			if err != nil {
				log.Fatal(err)
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

				for i, e := range v.testcase.Expectations {
					if e.Fulfilled || e.Pattern != pattern {
						continue
					}

					tokens := v.tokenizer.Tokenize(line, v.config.Patterns)

					// handle already verified expectations
					if e.Verified > 0 && e.Equal(tokens) {
						log.Printf("expectation verified by: %s\n", domain.Expectation{Tokens: tokens}.Shorten(6))
						v.testcase.Expectations[i].Fulfilled = true
						v.testcase.Expectations[i].Verified = e.Verified + 1
						break
					}

					// handle not yet verified expectations (verified == 0)
					if diff, err := e.Diff(tokens); err == nil {
						log.Printf("reference expectation found: %s\n", domain.Expectation{Tokens: tokens}.Shorten(6))
						v.testcase.Expectations[i].IgnoreDiffs = diff
						v.testcase.Expectations[i].Fulfilled = true
						v.testcase.Expectations[i].Verified = 1
						break
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
func (v *Verifier) ReportResults() domain.Report {
	fulfilled := 0
	verifiedSum := 0
	for _, e := range v.testcase.Expectations {
		verifiedSum += e.Verified
		if e.Fulfilled {
			fulfilled = fulfilled + 1
		}
	}
	report := domain.Report{
		Testname:         v.name,
		LastExecution:    time.Now(),
		Expectations:     len(v.testcase.Expectations),
		Verifications:    v.testcase.Verifications,
		Fulfilled:        fulfilled,
		VerificationMean: verificationMean(float32(verifiedSum), float32(len(v.testcase.Expectations))),
	}
	for _, e := range v.testcase.Expectations {
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
