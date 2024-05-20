package verify

import (
	"github.com/rwirdemann/datafrog/internal/datafrog"
	"github.com/rwirdemann/datafrog/internal/datafrog/mysql"
	"testing"

	"github.com/rwirdemann/datafrog/adapter"
	"github.com/stretchr/testify/assert"
)

func TestVerify(t *testing.T) {
	var emptyDiff []int
	var emptyExpectations []datafrog.Expectation
	testCases := []struct {
		desc                         string
		initialExpectations          []datafrog.Expectation
		updatedExpectations          []datafrog.Expectation
		additionalExpectations       []datafrog.Expectation
		logs                         []string
		patterns                     []string
		reportAdditionalExpectations bool
	}{
		{
			desc: "empty diff vector",
			initialExpectations: []datafrog.Expectation{
				{
					Tokens:      datafrog.Tokenize("select * from jobs;"),
					Pattern:     "select *",
					IgnoreDiffs: emptyDiff,
				},
			},
			logs: []string{
				"2024-04-08T09:39:15.070009Z	 2549 Query	select * from jobs;",
				"STOP",
			},
			updatedExpectations: []datafrog.Expectation{
				{
					Tokens:      datafrog.Tokenize("select * from jobs;"),
					Pattern:     "select *",
					IgnoreDiffs: emptyDiff,
					Fulfilled:   true,
					Verified:    1,
				},
			},
			additionalExpectations:       emptyExpectations,
			patterns:                     []string{"select *"},
			reportAdditionalExpectations: true,
		},
		{
			desc: "diff vector with one element",
			initialExpectations: []datafrog.Expectation{
				{
					Tokens:      datafrog.Tokenize("select * from jobs where id=1;"),
					Pattern:     "select *",
					IgnoreDiffs: []int{},
				},
			},
			logs: []string{
				"2024-04-08T09:39:15.070009Z	 2549 Query	select * from jobs where id=2;",
				"STOP",
			},
			updatedExpectations: []datafrog.Expectation{
				{
					Tokens:      datafrog.Tokenize("select * from jobs where id=1;"),
					Pattern:     "select *",
					IgnoreDiffs: []int{5},
					Fulfilled:   true,
					Verified:    1,
				},
			},
			additionalExpectations:       emptyExpectations,
			patterns:                     []string{"select *"},
			reportAdditionalExpectations: true,
		},
		{
			desc: "miss matching pattern",
			initialExpectations: []datafrog.Expectation{
				{
					Tokens:      datafrog.Tokenize("select * from jobs where id=1;"),
					Pattern:     "select *",
					IgnoreDiffs: []int{},
				},
			},
			logs: []string{
				"2024-04-08T09:39:15.070009Z	 2549 Query	update job where id=2;",
				"STOP",
			},
			updatedExpectations: []datafrog.Expectation{
				{
					Tokens:      datafrog.Tokenize("select * from jobs where id=1;"),
					Pattern:     "select *",
					IgnoreDiffs: []int{},
					Fulfilled:   false,
					Verified:    0,
				},
			},
			additionalExpectations:       emptyExpectations,
			patterns:                     []string{"select *"},
			reportAdditionalExpectations: true,
		},
		{
			desc: "multiple expectations, different order",
			initialExpectations: []datafrog.Expectation{
				{
					Tokens:      datafrog.Tokenize("select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0"),
					Pattern:     "select",
					IgnoreDiffs: []int{},
				},
				{
					Tokens:      datafrog.Tokenize("insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World', '2024-04-17 15:55:56', 0, null, '', 'Hello', 3)"),
					Pattern:     "insert",
					IgnoreDiffs: []int{},
				},
				{
					Tokens:      datafrog.Tokenize("update job set description='World, X', publish_at='2024-04-17 15:55:56', publish_trials=1, published_timestamp='2024-04-17 15:55:58.433346', tags='', title='Hello' where id=3"),
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
			updatedExpectations: []datafrog.Expectation{
				{
					Tokens:      datafrog.Tokenize("select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_"),
					Pattern:     "select",
					IgnoreDiffs: emptyDiff,
					Fulfilled:   true,
					Verified:    1,
				},
				{
					Tokens:      datafrog.Tokenize("insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World', '2024-04-17 15:55:56', 0, null, '', 'Hello', 4)"),
					Pattern:     "insert",
					IgnoreDiffs: []int{12, 17},
					Fulfilled:   true,
					Verified:    1,
				},
				{
					Tokens:      datafrog.Tokenize("update job set description='World, X', publish_at='2024-04-17 15:55:56', publish_trials=1, published_timestamp='2024-04-17 15:55:58.433346', tags='', title='Hello' where id=3"),
					Pattern:     "update",
					IgnoreDiffs: []int{6, 10},
					Fulfilled:   true,
					Verified:    1,
				},
			},
			additionalExpectations:       emptyExpectations,
			patterns:                     []string{"select", "insert", "update"},
			reportAdditionalExpectations: true,
		},
		{
			desc: "verified > 0 but equals fails should continue with next expectation",
			initialExpectations: []datafrog.Expectation{
				{
					Tokens:      datafrog.Tokenize("insert into job (description, id) values ('Developer', 4)"),
					Pattern:     "insert",
					Verified:    1,
					Fulfilled:   true,
					IgnoreDiffs: []int{7},
				},
				{
					Tokens:      datafrog.Tokenize("insert into application (name, job_id, id) values ('Ralf', 4, 1)"),
					Pattern:     "insert",
					Verified:    1,
					Fulfilled:   true,
					IgnoreDiffs: []int{8, 9},
				},
			},
			logs: []string{
				"2024-04-08T09:39:15.070009Z	 2549 Query	insert into application (name, job_id, id) values ('Ralf', 5, 2);",
				"2024-04-08T09:39:15.070009Z	 2549 Query	insert into job (description, id) values ('Developer', 5);",
				"STOP",
			},
			updatedExpectations: []datafrog.Expectation{
				{
					Tokens:      datafrog.Tokenize("insert into job (description, id) values ('Developer', 4)"),
					Pattern:     "insert",
					Verified:    2,
					Fulfilled:   true,
					IgnoreDiffs: []int{7},
				},
				{
					Tokens:      datafrog.Tokenize("insert into application (name, job_id, id) values ('Ralf', 4, 1)"),
					Pattern:     "insert",
					Verified:    2,
					Fulfilled:   true,
					IgnoreDiffs: []int{8, 9},
				},
			},
			additionalExpectations:       emptyExpectations,
			patterns:                     []string{"insert"},
			reportAdditionalExpectations: true,
		},
		{
			desc: "additional pattern matching verifications",
			initialExpectations: []datafrog.Expectation{
				{
					Tokens:      datafrog.Tokenize("select * from jobs;"),
					Pattern:     "select *",
					IgnoreDiffs: emptyDiff,
				},
			},
			logs: []string{
				"2024-04-08T09:39:15.070009Z	 2549 Query	insert into jobs;",
				"2024-04-08T09:39:16.070009Z	 2550 Query	select * from jobs;",
				"STOP",
			},
			updatedExpectations: []datafrog.Expectation{
				{
					Tokens:      datafrog.Tokenize("select * from jobs;"),
					Pattern:     "select *",
					IgnoreDiffs: emptyDiff,
					Fulfilled:   true,
					Verified:    1,
				},
			},
			additionalExpectations: []datafrog.Expectation{
				{
					Tokens:      datafrog.Tokenize("insert into jobs;"),
					Pattern:     "insert into",
					IgnoreDiffs: emptyDiff,
					Fulfilled:   false,
					Verified:    0,
				},
			},
			patterns:                     []string{"select *", "insert into"},
			reportAdditionalExpectations: true,
		},
		{
			desc: "additional expecatations should not be reported",
			initialExpectations: []datafrog.Expectation{
				{
					Tokens:      datafrog.Tokenize("select * from jobs;"),
					Pattern:     "select *",
					IgnoreDiffs: emptyDiff,
				},
			},
			logs: []string{
				"2024-04-08T09:39:15.070009Z	 2549 Query	insert into jobs;",
				"2024-04-08T09:39:16.070009Z	 2550 Query	select * from jobs;",
				"STOP",
			},
			updatedExpectations: []datafrog.Expectation{
				{
					Tokens:      datafrog.Tokenize("select * from jobs;"),
					Pattern:     "select *",
					IgnoreDiffs: emptyDiff,
					Fulfilled:   true,
					Verified:    1,
				},
			},
			additionalExpectations:       emptyExpectations,
			patterns:                     []string{"select *", "insert into"},
			reportAdditionalExpectations: false,
		},
		{
			// This test reproduces the following error situation:
			// - we have an already verified expectation e1 (verified > 0)
			// - e1 matches the pattern of the current v but isn't equal to v
			// - e1 anbd v have the same token length
			// => Bug: a new diff vector was build and stored for e1
			desc: "fix diff building for already verfied expecation",
			initialExpectations: []datafrog.Expectation{
				{
					Tokens:      datafrog.Tokenize("insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World', '2024-04-17 15:55:56', 0, null, '', 'Hello', 3)"),
					Pattern:     "insert",
					Verified:    1,
					Fulfilled:   true,
					IgnoreDiffs: []int{12, 17},
				},
			},
			logs: []string{
				"2024-04-17T13:55:56.750174Z\t 2000 Query\tinsert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('Universe', '2024-04-17 15:56:56', 0, null, '', 'Yeah', 4)",
				"STOP",
			},
			updatedExpectations: []datafrog.Expectation{
				{
					Tokens:      datafrog.Tokenize("insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World', '2024-04-17 15:55:56', 0, null, '', 'Hello', 4)"),
					Pattern:     "insert",
					IgnoreDiffs: []int{12, 17},
					Fulfilled:   false,
					Verified:    1,
				},
			},
			additionalExpectations:       emptyExpectations,
			patterns:                     []string{"insert"},
			reportAdditionalExpectations: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			c := datafrog.Config{}
			c.Patterns = tC.patterns
			c.Expectations.ReportAdditional = tC.reportAdditionalExpectations
			doneChannel := make(chan struct{})
			stoppedChannel := make(chan struct{})
			databaseLog := adapter.NewMemSQLLog(tC.logs, doneChannel)
			expectationSource := adapter.NewMemExpectationSource(tC.initialExpectations)
			timer := adapter.MockTimer{}
			verifier := NewVerifier(c, mysql.Tokenizer{}, databaseLog, expectationSource, timer, "")
			go verifier.Start(doneChannel, stoppedChannel)
			<-stoppedChannel // wait till verifier is done
			tc := expectationSource.Get()
			for i, e := range tc.Expectations {
				updatedExpectation := tC.updatedExpectations[i]
				assert.Equal(t, updatedExpectation.IgnoreDiffs, e.IgnoreDiffs)
				assert.Equal(t, updatedExpectation.Fulfilled, e.Fulfilled)
				assert.Equal(t, updatedExpectation.Verified, e.Verified)
			}
			assert.Equal(t, tC.additionalExpectations, tc.AdditionalExpectations)
		})
	}
}
