package cmd

import (
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
	}{
		{
			desc: "",
			expectations: []string{
				"2024-04-08T09:36:15.070009Z	 2549 Query	select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ order by job0_.publish_at desc",
			},
			logs: []string{
				"2024-04-08T09:39:15.070009Z	 2549 Query	select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ order by job0_.publish_at desc",
				"STOP",
			},
			patterns:              []string{"insert", "select job!publish_trials<1"},
			remainingExpectations: 0,
			recordedVerifications: 1,
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
			verifier.Start()
			verifier.Stop()
			assert.Len(t, expectationSource.GetAll(), tC.remainingExpectations)
			assert.Len(t, verificationSink.Recorded, tC.recordedVerifications)
		})
	}
}
