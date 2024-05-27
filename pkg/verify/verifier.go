package verify

import (
	"encoding/json"
	"github.com/rwirdemann/datafrog/pkg/df"
	"log"
	"time"
)

// The Verifier verifies the expectations in expectationSource. It monitors the
// log for these expectations and increases their verify count if matched. The
// updated expectation list is written back to expectationSource after the
// verification run is done.
type Verifier struct {
	config    df.Config
	tokenizer df.Tokenizer
	log       df.Log
	writer    df.TestWriter // to write the verfied and updated testscase
	testcase  df.Testcase
	timer     df.Timer
	name      string
}

// NewVerifier creates a new Verifier.
func NewVerifier(
	c df.Config,
	tokenizer df.Tokenizer,
	log df.Log,
	tc df.Testcase,
	w df.TestWriter,
	t df.Timer,
	name string) *Verifier {
	return &Verifier{
		config:    c,
		tokenizer: tokenizer,
		log:       log,
		writer:    w,
		testcase:  tc,
		timer:     t,
		name:      name,
	}
}
func (verifier *Verifier) Testcase() df.Testcase {
	return verifier.testcase
}

// Start runs the verification loop. Stops when done channel was closed. Closes
// stopped channel afterward in order to tell its caller (web, cli, ...) that
// verification has been finished.
func (verifier *Verifier) Start(done chan struct{}, stopped chan struct{}) {
	verifier.timer.Start()
	log.Printf("verification started at %v...", verifier.timer.GetStart())
	verifier.testcase.Verifications = verifier.testcase.Verifications + 1
	verifier.testcase.LastExecution = time.Now()
	for i := range verifier.testcase.Expectations {
		verifier.testcase.Expectations[i].Fulfilled = false
	}

	// tell caller that verification has been finished
	defer close(stopped)

	// called when done channel is closed
	defer func() {
		verifier.write()
	}()

	// jump to log file end
	if err := verifier.log.Tail(); err != nil {
		log.Fatal(err)
	}

	for {
		select {
		default:
			v, err := verifier.log.NextLine()
			if err != nil {
				log.Fatal(err)
			}

			ts, err := verifier.log.Timestamp(v)
			if err != nil {
				continue
			}
			if verifier.timer.MatchesRecordingPeriod(ts) {
				matches, vPattern := df.MatchesPattern(verifier.config, v)
				if !matches {
					continue
				}

				verified := verifier.verify(v, vPattern)

				if !verified && verifier.config.Expectations.ReportAdditional {

					// v matches pattern but no matching expectation was found
					expectation := df.Expectation{
						Tokens: verifier.tokenizer.Tokenize(v, verifier.config.Patterns), Pattern: vPattern,
					}
					log.Printf("additional expectation found: %s\n", expectation.Shorten(6))
					verifier.testcase.AdditionalExpectations = append(verifier.testcase.AdditionalExpectations, expectation)
				}
			}
		case <-done:
			log.Printf("verifier: done channel closed")
			return
		}
	}
}

func (verifier *Verifier) write() {
	// create a write copy of the testcase to make sure no addtional
	// expectations are saved but kept for reporting reasons
	tc := verifier.testcase

	// don't write additional expectations
	tc.AdditionalExpectations = nil

	b, err := json.Marshal(tc)
	if err != nil {
		log.Fatal(err)
	}

	if len(b) == 0 {
		log.Fatal("testcase is empty. skip writing")
	}

	_, err = verifier.writer.Write(b)
	if err != nil {
		log.Fatal(err)
	}
	if err := verifier.writer.Close(); err != nil {
		log.Fatal(err)
	}
	log.Printf("verifier: testcase successfully written back")
}

// verify tries to verify one of the testcases expectations. Returns true if an
// expectation was verfied and false otherwise.
func (verifier *Verifier) verify(v string, vPattern string) bool {
	for i, e := range verifier.testcase.Expectations {
		if e.Fulfilled || e.Pattern != vPattern {
			continue // -> continue with next e
		}

		vTokens := verifier.tokenizer.Tokenize(v, verifier.config.Patterns)

		// Handle already verified expectations (reference expectation)
		if e.Verified > 0 && e.Equal(vTokens) {
			log.Printf("expectation verified by: %s\n", df.Expectation{Tokens: vTokens}.Shorten(6))
			verifier.testcase.Expectations[i].Fulfilled = true
			verifier.testcase.Expectations[i].Verified = e.Verified + 1
			return true // -> continue with next v
		}

		if len(e.Tokens) != len(vTokens) {
			continue // -> continue with next e
		}

		if e.Verified == 0 {
			// Not yet verified expectation e with same token lengths as v
			// found. This expectation e becomes our reference expectation.
			if diff, err := e.Diff(vTokens); err == nil {
				log.Printf("reference expectation found: %s\n", df.Expectation{Tokens: vTokens}.Shorten(6))
				verifier.testcase.Expectations[i].IgnoreDiffs = diff
				verifier.testcase.Expectations[i].Fulfilled = true
				verifier.testcase.Expectations[i].Verified = 1
				return true // -> continue with next v
			}
		}
	}
	return false // -> expectation not verified
}

// ReportResults creates a [domain.Report] of the verification results.
func (verifier *Verifier) ReportResults() df.Report {
	fulfilled := 0
	verifiedSum := 0
	for _, e := range verifier.testcase.Expectations {
		verifiedSum += e.Verified
		if e.Fulfilled {
			fulfilled = fulfilled + 1
		}
	}
	report := df.Report{
		Testname:         verifier.name,
		LastExecution:    time.Now(),
		Expectations:     len(verifier.testcase.Expectations),
		Verifications:    verifier.testcase.Verifications,
		Fulfilled:        fulfilled,
		VerificationMean: verificationMean(float32(verifiedSum), float32(len(verifier.testcase.Expectations))),
	}
	for _, e := range verifier.testcase.Expectations {
		if !e.Fulfilled {
			report.Unfulfilled = append(report.Unfulfilled, e)
		}
	}
	for _, e := range verifier.testcase.AdditionalExpectations {
		report.AdditionalExpectations = append(report.AdditionalExpectations, e.Shorten(6))
	}
	return report
}

func verificationMean(sum, expectationCount float32) float32 {
	if expectationCount > 0 {
		return sum / expectationCount
	}
	return 0
}
