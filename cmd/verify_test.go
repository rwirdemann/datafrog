package cmd

import (
	"testing"

	"github.com/rwirdemann/databasedragon/adapter"
	"github.com/rwirdemann/databasedragon/config"
	"github.com/rwirdemann/databasedragon/matcher"
	"github.com/stretchr/testify/assert"
)

func TestVerify(t *testing.T) {
	testCases := []struct {
		desc         string
		expectations []matcher.Expectation
		logs         []string
		patterns     []string
		ignoreDiffs  []int
	}{
		// {
		// 	desc: "success",
		// 	expectations: []matcher.Expectation{
		// 		{
		// 			Tokens:      matcher.Tokenize("select * from jobs;"),
		// 			Pattern:     "select job",
		// 			IgnoreDiffs: []int{},
		// 		},
		// 	},
		// 	logs: []string{
		// 		"2024-04-08T09:39:15.070009Z	 2549 Query	select * from jobs;",
		// 		"STOP",
		// 	},
		// 	patterns:    []string{"select job"},
		// 	ignoreDiffs: []int{},
		// },
		{
			desc: "with id",
			expectations: []matcher.Expectation{
				{
					Tokens:      matcher.Tokenize("select * from jobs where id=1;"),
					Pattern:     "select *",
					IgnoreDiffs: []int{},
				},
			},
			logs: []string{
				"2024-04-08T09:39:15.070009Z	 2549 Query	select * from jobs where id=2;",
				"STOP",
			},
			patterns:    []string{"select *"},
			ignoreDiffs: []int{5},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			c := config.Config{}
			c.Patterns = tC.patterns
			databaseLog := adapter.NewMemSQLLog(tC.logs)
			expectationSource := adapter.NewMemExpectationSource(tC.expectations)
			timer := adapter.MockTimer{}
			verifier := NewVerifier(c, databaseLog, expectationSource, timer)
			verifier.Start()
			updatedExpectations := expectationSource.GetAll()
			assert.Equal(t, tC.ignoreDiffs, updatedExpectations[0].IgnoreDiffs)
		})
	}
}
