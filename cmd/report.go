package cmd

import (
	"fmt"
	"strings"
	"time"
)

type Report struct {
	Testname      string    `json:"testname"`
	LastExecution time.Time `json:"last_execution"`
	Expectations  int       `json:"expectations"`
	Fulfilled     int       `json:"fulfilled"`
	Unfulfilled   []string  `json:"unfulfilled"`
}

func (r Report) String() string {
	return fmt.Sprintf("Testname: %s\n"+
		"Last execution: %s\n"+
		"Expections: %d\n"+
		"Fulfilled: %d\n"+
		"Unfulfilled: %s\n",
		r.Testname,
		r.LastExecution.Format(time.DateTime),
		r.Expectations,
		r.Fulfilled,
		strings.Join(r.Unfulfilled, "\n"))
}
