package df

import (
	"fmt"
	"strings"
	"time"
)

type Report struct {
	Testname               string        `json:"testname"`
	LastExecution          time.Time     `json:"last_execution"`
	Verifications          int           `json:"verifications"`
	Expectations           int           `json:"expectations"`
	Fulfilled              int           `json:"fulfilled"`
	Unfulfilled            []Expectation `json:"unfulfilled,omitempty"`
	VerificationMean       float32       `json:"verification_mean"`
	AdditionalExpectations []string      `json:"additional_expectations,omitempty"`
}

func (r Report) String() string {

	return fmt.Sprintf("Testname: %s\n"+
		"Last execution: %s\n"+
		"Verfications: %d\n"+
		"Expectations: %d\n"+
		"Fulfilled: %d\n"+
		"Verification mean: %f\n"+
		"Unfulfilled: %s\n",
		r.Testname,
		r.LastExecution.Format(time.DateTime),
		r.Verifications,
		r.Expectations,
		r.Fulfilled,
		r.VerificationMean,
		strings.Join(toString(r.Unfulfilled), "\n"))
}

func toString(e []Expectation) []string {
	var result []string
	for _, e := range e {
		result = append(result, e.Shorten(6))
	}
	return result
}
