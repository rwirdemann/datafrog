package domain

import (
	"fmt"
	"strings"
	"time"
)

type Report struct {
	Testname         string    `json:"testname"`
	LastExecution    time.Time `json:"last_execution"`
	Expectations     int       `json:"expectations"`
	Fulfilled        int       `json:"fulfilled"`
	Unfulfilled      []string  `json:"unfulfilled,omitempty"`
	MaxVerified      int       `json:"max_verified"`
	VerificationMean float32   `json:"verification_mean"`
}

func (r Report) String() string {
	return fmt.Sprintf("Testname: %s\n"+
		"Last execution: %s\n"+
		"Expectations: %d\n"+
		"Fulfilled: %d\n"+
		"Verification mean: %f\n"+
		"Unfulfilled: %s\n",
		r.Testname,
		r.LastExecution.Format(time.DateTime),
		r.Expectations,
		r.Fulfilled,
		r.VerificationMean,
		strings.Join(r.Unfulfilled, "\n"))
}
