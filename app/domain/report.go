package domain

import (
	"fmt"
	"strings"
	"time"
)

type Report struct {
	Testname         string    `json:"testname"`
	LastExecution    time.Time `json:"last_execution"`
	Verifications    int       `json:"verifications"`
	Expectations     int       `json:"expectations"`
	Fulfilled        int       `json:"fulfilled"`
	Unfulfilled      []string  `json:"unfulfilled,omitempty"`
	VerificationMean float32   `json:"verification_mean"`
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
		strings.Join(r.Unfulfilled, "\n"))
}
