package record

import (
	"github.com/rwirdemann/datafrog/pkg/df"
	"github.com/rwirdemann/datafrog/pkg/mocks"
	"github.com/rwirdemann/datafrog/pkg/mysql"
	"testing"

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

	e1 := df.Expectation{
		Tokens:      df.Tokenize("select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ order by job0_.publish_at desc"),
		IgnoreDiffs: []int{},
		Verified:    0,
		Fulfilled:   false,
		Pattern:     "select job!publish_trials<1",
		Uuid:        "023a6a95-6c8a-4483-bcfb-17b1c58c317f",
	}

	e2 := df.Expectation{
		Tokens:      df.Tokenize("insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World', '2024-04-08 14:50:20', 0, null, '', 'Hello', 3)"),
		IgnoreDiffs: []int{},
		Verified:    0,
		Fulfilled:   false,
		Pattern:     "insert",
		Uuid:        "023a6a95-6c8a-4483-bcfb-17b1c58c317f",
	}

	expectedTestcase := df.Testcase{
		Name:          "create-job.json",
		Running:       false,
		Verifications: 0,
		Expectations:  []df.Expectation{e1, e2},
	}

	c := df.Config{}
	c.Patterns = []string{"insert", "select job!publish_trials<1"}

	recordingDone := make(chan struct{})
	recordingStopped := make(chan struct{})
	databaseLog := mocks.NewMemSQLLog(logs, recordingDone)
	writer := &mocks.MemWriter{}
	timer := mocks.Timer{}
	recorder := NewRecorder(c, mysql.Tokenizer{}, databaseLog, writer, timer, "create-job.json", mocks.StaticUUIDProvider{})
	go recorder.Start(recordingDone, recordingStopped)
	<-recordingStopped
	assert.Len(t, writer.Testcase.Expectations, 2)
	assert.Equal(t, expectedTestcase, writer.Testcase)
}
