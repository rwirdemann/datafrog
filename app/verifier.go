package app

import (
	"fmt"
	"github.com/rwirdemann/datafrog/app/domain"
	"github.com/rwirdemann/datafrog/config"
	"github.com/rwirdemann/datafrog/matcher"
	"github.com/rwirdemann/datafrog/ports"
	"log"
	"time"
)

// The Verifier verifies the expectations in expectationSource. It monitors the
// log for these expectations and increases their verify count if matched. The
// updated expectation list is written back to expectationSource after the
// verification run is done.
type Verifier struct {
	config            config.Config
	tokenizer         matcher.Tokenizer
	log               ports.Log
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
		log:               log,
		expectationSource: source,
		testcase:          source.Get(),
		timer:             t,
		name:              name,
	}
}

// Start runs the verification loop. Stops when done channel was closed. Closes
// stopped channel afterward in order to tell its caller (web, cli, ...) that
// verification has been finished.
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

			v, err := verifier.log.NextLine()
			if err != nil {
				log.Fatal(err)
			}

			ts, err := verifier.log.Timestamp(v)
			if err != nil {
				continue
			}
			if verifier.timer.MatchesRecordingPeriod(ts) {
				matches, vPattern := matcher.MatchesPattern(verifier.config, v)
				if !matches {
					continue
				}

				expectationVerified := false
				for i, e := range verifier.testcase.Expectations {
					if e.Fulfilled || e.Pattern != vPattern {
						continue // -> continue with next e
					}

					vTokens := verifier.tokenizer.Tokenize(v, verifier.config.Patterns)

					// Handle already verified expectations.
					if e.Verified > 0 && e.Equal(vTokens) {
						log.Printf("expectation verified by: %s\n", domain.Expectation{Tokens: vTokens}.Shorten(6))
						verifier.testcase.Expectations[i].Fulfilled = true
						verifier.testcase.Expectations[i].Verified = e.Verified + 1
						expectationVerified = true
						break // -> continue with next v
					}

					if len(e.Tokens) != len(vTokens) {
						continue // -> continue with next e
					}

					// Not yet verified expectation e with same token lengths as
					// v found. This expectation e becomes our references
					// expectation.
					if diff, err := e.Diff(vTokens); err == nil {
						log.Printf("reference expectation found: %s\n", domain.Expectation{Tokens: vTokens}.Shorten(6))
						verifier.testcase.Expectations[i].IgnoreDiffs = diff
						verifier.testcase.Expectations[i].Fulfilled = true
						verifier.testcase.Expectations[i].Verified = 1
						expectationVerified = true
						break // -> continue with next v
					}
				}

				// v matches pattern but no matching expectation was found
				if !expectationVerified {
					expectation := domain.Expectation{
						Tokens: verifier.tokenizer.Tokenize(v, verifier.config.Patterns), Pattern: vPattern,
					}
					log.Printf("additional expectation found: %s\n", expectation.Shorten(6))
					verifier.testcase.AdditionalExpectations = append(verifier.testcase.AdditionalExpectations, expectation)
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

// ReportResults creates a [domain.Report] of the verification results.
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
