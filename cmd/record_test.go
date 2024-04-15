package cmd

import (
	"encoding/json"
	"testing"

	"github.com/rwirdemann/databasedragon/adapter"
	"github.com/rwirdemann/databasedragon/config"
	"github.com/rwirdemann/databasedragon/matcher"
	"github.com/stretchr/testify/assert"
)

func TestRecord(t *testing.T) {
	logs := []string{
		"2024-04-08T09:33:15.070009Z	 2549 Query	select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ order by job0_.publish_at desc",
		"2024-04-08T12:50:59.605638Z	 2609 Query	insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World', '2024-04-08 14:50:20', 0, null, '', 'Hello', 3)",

		// do not record due to unmatched or excluded patterns
		"2024-04-08T12:47:14.012398Z	 2609 Query	update job set description='World', publish_at='2024-04-08 14:47:04', publish_trials=1, published_timestamp='2024-04-08 14:47:14.006028', tags='', title='Hello' where id=1",
		"2024-04-08T15:40:58.414756Z	 2669 Query	select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ where (job0_.published_timestamp is null) and job0_.publish_at<'2024-04-08 17:40:58.414289' and job0_.publish_trials<1",

		"STOP",
	}

	e1 := matcher.Expectation{
		Tokens:      matcher.Tokenize("select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ order by job0_.publish_at desc"),
		IgnoreDiffs: []int{},
	}

	e2 := matcher.Expectation{
		Tokens:      matcher.Tokenize("insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World', '2024-04-08 14:50:20', 0, null, '', 'Hello', 3)"),
		IgnoreDiffs: []int{},
	}

	expectations, _ := json.Marshal([]matcher.Expectation{e1, e2})

	c := config.Config{}
	c.Patterns = []string{"insert", "select job!publish_trials<1"}

	databaseLog := adapter.NewMemSQLLog(logs)
	recordingSink := adapter.NewMemRecordingSink()
	timer := adapter.MockTimer{}
	recorder := NewRecorder(c, databaseLog, recordingSink, timer)
	recorder.Start()
	recorder.Stop()
	assert.Len(t, recordingSink.Recorded, 1)
	assert.Equal(t, string(expectations), recordingSink.Recorded[0])
}
