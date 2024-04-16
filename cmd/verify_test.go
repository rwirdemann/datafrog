package cmd

import (
	"testing"

	"github.com/rwirdemann/databasedragon/adapter"
	"github.com/rwirdemann/databasedragon/config"
	"github.com/rwirdemann/databasedragon/matcher"
	"github.com/stretchr/testify/assert"
)

func TestVerify(t *testing.T) {
	var emptyDiff []int
	testCases := []struct {
		desc                string
		initialExpectations []matcher.Expectation
		updatedExpectations []matcher.Expectation
		logs                []string
		patterns            []string
	}{
		{
			desc: "empty diff vector",
			initialExpectations: []matcher.Expectation{
				{
					Tokens:      matcher.Tokenize("select * from jobs;"),
					Pattern:     "select *",
					IgnoreDiffs: emptyDiff,
				},
			},
			logs: []string{
				"2024-04-08T09:39:15.070009Z	 2549 Query	select * from jobs;",
				"STOP",
			},
			updatedExpectations: []matcher.Expectation{
				{
					Tokens:      matcher.Tokenize("select * from jobs;"),
					Pattern:     "select *",
					IgnoreDiffs: emptyDiff,
					Fulfilled:   true,
				},
			},
			patterns: []string{"select *"},
		},
		{
			desc: "diff vector with one element",
			initialExpectations: []matcher.Expectation{
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
			updatedExpectations: []matcher.Expectation{
				{
					Tokens:      matcher.Tokenize("select * from jobs where id=1;"),
					Pattern:     "select *",
					IgnoreDiffs: []int{5},
					Fulfilled:   true,
				},
			},
			patterns: []string{"select *"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			c := config.Config{}
			c.Patterns = tC.patterns
			databaseLog := adapter.NewMemSQLLog(tC.logs)
			expectationSource := adapter.NewMemExpectationSource(tC.initialExpectations)
			timer := adapter.MockTimer{}
			verifier := NewVerifier(c, databaseLog, expectationSource, timer)
			verifier.Start()
			expectations := expectationSource.GetAll()
			for i, e := range expectations {
				updatedExpectation := tC.updatedExpectations[i]
				assert.Equal(t, updatedExpectation.IgnoreDiffs, e.IgnoreDiffs)
				assert.Equal(t, updatedExpectation.Fulfilled, e.Fulfilled)
			}

		})
	}
}
