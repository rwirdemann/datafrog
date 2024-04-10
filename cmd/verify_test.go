package cmd

import (
	"errors"
	"testing"

	"github.com/rwirdemann/databasedragon/adapter"
	"github.com/rwirdemann/databasedragon/config"
	"github.com/stretchr/testify/assert"
)

func TestVerify(t *testing.T) {
	testCases := []struct {
		desc                  string
		expectations          []string
		logs                  []string
		patterns              []string
		remainingExpectations int
		recordedVerifications int
		expectedError         error
	}{
		{
			desc: "success",
			expectations: []string{
				"2024-04-08T09:36:15.070009Z	 2549 Query	select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ order by job0_.publish_at desc",
			},
			logs: []string{
				"2024-04-08T09:39:15.070009Z	 2549 Query	select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ order by job0_.publish_at desc",
				"STOP",
			},
			patterns:              []string{"select job!publish_trials<1"},
			remainingExpectations: 0,
			recordedVerifications: 1,
			expectedError:         nil,
		},
		{
			desc: "first log matches pattern but not first expectation",
			expectations: []string{
				"2024-04-08T09:36:15.070009Z	 2549 Query	select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ order by job0_.publish_at desc",
			},
			logs: []string{
				"2024-04-02T06:38:05.015501Z     1669 Query	update job set description='World, X', publish_at='2024-04-02 08:37:37', tags='', title='Hello' where id=39",
				"2024-04-08T09:39:15.070009Z	 2549 Query	select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ order by job0_.publish_at desc",
				"STOP",
			},
			patterns:              []string{"select job!publish_trials<1", "update"},
			remainingExpectations: 1,
			recordedVerifications: 0,
			expectedError:         errors.New("first expectation didn't match expected pattern"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			c := config.Config{}
			c.Patterns = tC.patterns
			databaseLog := adapter.NewMemSQLLog(tC.logs)
			expectationSource := adapter.NewMemExpectationSource(tC.expectations)
			verificationSink := adapter.NewMemRecordingSink()
			timer := adapter.MockTimer{}
			verifier := NewVerifier(c, databaseLog, expectationSource, verificationSink, timer)
			err := verifier.Start()
			assert.Equal(t, tC.expectedError, err)
			verifier.Stop()
			assert.Len(t, expectationSource.GetAll(), tC.remainingExpectations)
			assert.Len(t, verificationSink.Recorded, tC.recordedVerifications)
		})
	}
}
