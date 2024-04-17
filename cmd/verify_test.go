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
					Verified:    1,
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
					Verified:    1,
				},
			},
			patterns: []string{"select *"},
		},
		{
			desc: "miss matching pattern",
			initialExpectations: []matcher.Expectation{
				{
					Tokens:      matcher.Tokenize("select * from jobs where id=1;"),
					Pattern:     "select *",
					IgnoreDiffs: []int{},
				},
			},
			logs: []string{
				"2024-04-08T09:39:15.070009Z	 2549 Query	update job where id=2;",
				"STOP",
			},
			updatedExpectations: []matcher.Expectation{
				{
					Tokens:      matcher.Tokenize("select * from jobs where id=1;"),
					Pattern:     "select *",
					IgnoreDiffs: []int{},
					Fulfilled:   false,
					Verified:    0,
				},
			},
			patterns: []string{"select *"},
		},
		{
			desc: "multiple expectations, different order",
			initialExpectations: []matcher.Expectation{
				{
					Tokens:      matcher.Tokenize("select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0"),
					Pattern:     "select",
					IgnoreDiffs: []int{},
				},
				{
					Tokens:      matcher.Tokenize("insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World', '2024-04-17 15:55:56', 0, null, '', 'Hello', 3)"),
					Pattern:     "insert",
					IgnoreDiffs: []int{},
				},
				{
					Tokens:      matcher.Tokenize("update job set description='World, X', publish_at='2024-04-17 15:55:56', publish_trials=1, published_timestamp='2024-04-17 15:55:58.433346', tags='', title='Hello' where id=3"),
					Pattern:     "update",
					IgnoreDiffs: []int{},
				},
			},
			logs: []string{
				"2024-04-17T13:55:56.750174Z\t 2000 Query\tinsert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World', '2024-04-17 15:56:56', 0, null, '', 'Hello', 4)",
				"2024-04-17T13:55:58.434784Z\t 2000 Query\tupdate job set description='World, X', publish_at='2024-04-17 15:55:56', publish_trials=1, published_timestamp='2024-04-17 15:55:59.433346', tags='', title='Hello' where id=4",
				"2024-04-17T13:55:57.090960Z\t 2001 Query\tselect job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0",
				"STOP",
			},
			updatedExpectations: []matcher.Expectation{
				{
					Tokens:      matcher.Tokenize("select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_"),
					Pattern:     "select",
					IgnoreDiffs: emptyDiff,
					Fulfilled:   true,
					Verified:    1,
				},
				{
					Tokens:      matcher.Tokenize("insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World', '2024-04-17 15:55:56', 0, null, '', 'Hello', 4)"),
					Pattern:     "insert",
					IgnoreDiffs: []int{12, 17},
					Fulfilled:   true,
					Verified:    1,
				},
				{
					Tokens:      matcher.Tokenize("update job set description='World, X', publish_at='2024-04-17 15:55:56', publish_trials=1, published_timestamp='2024-04-17 15:55:58.433346', tags='', title='Hello' where id=3"),
					Pattern:     "update",
					IgnoreDiffs: []int{6, 10},
					Fulfilled:   true,
					Verified:    1,
				},
			},
			patterns: []string{"select", "insert", "update"},
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
				assert.Equal(t, updatedExpectation.Verified, e.Verified)
			}
		})
	}
}
