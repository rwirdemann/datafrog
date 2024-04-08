package cmd

import (
	"testing"

	"github.com/rwirdemann/databasedragon/adapter"
	"github.com/rwirdemann/databasedragon/config"
	"github.com/stretchr/testify/assert"
)

func TestListen(t *testing.T) {
	c := config.Config{}
	c.Patterns = []string{"select"}

	expectations := []string{
		"2024-04-08T09:36:15.070009Z	 2549 Query	select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ order by job0_.publish_at desc",
	}
	verifications := []string{
		"2024-04-08T09:36:13.070009Z	 2549 Query	select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ order by job0_.publish_at desc",
	}
	expectationSource := adapter.NewMemExpectationSource(expectations)
	verificationSource := adapter.NewMemExpectationSource(verifications)

	logs := []string{
		"2024-04-08T09:33:15.070009Z	 2549 Query	select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ order by job0_.publish_at desc",
		"STOP",
	}

	databaseLog := adapter.NewMemSQLLog(logs)
	timer := adapter.MockTimer{}
	listener = NewListener(c, timer, databaseLog, expectationSource, verificationSource)
	listener.Start()
	listener.Stop()
	results := listener.GetResults()
	assert.Len(t, results, 0)
}
