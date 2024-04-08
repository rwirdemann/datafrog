package cmd

import (
	"testing"

	"github.com/rwirdemann/databasedragon/adapter"
	"github.com/rwirdemann/databasedragon/config"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestListen(t *testing.T) {
	c := config.Config{}
	c.Patterns = []string{"select", "insert", "update"}

	expectations := []string{
		"2024-04-08T09:36:15.070009Z	 2549 Query	select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ order by job0_.publish_at desc",
		"2024-04-08T12:47:11.949960Z	 2609 Query	insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World', '2024-04-08 14:47:04', 0, null, '', 'Hello', 1)",
		"2024-04-08T12:47:11.953955Z	 2609 Query	select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ order by job0_.publish_at desc",
		"2024-04-08T12:47:14.012398Z	 2609 Query	update job set description='World', publish_at='2024-04-08 14:47:04', publish_trials=1, published_timestamp='2024-04-08 14:47:14.006028', tags='', title='Hello' where id=1",
		"2024-04-08T14:10:36.824558Z	 2639 Query	select job0_.id as id1_0_0_, job0_.description as descript2_0_0_, job0_.publish_at as publish_3_0_0_, job0_.publish_trials as publish_4_0_0_, job0_.published_timestamp as publishe5_0_0_, job0_.tags as tags6_0_0_, job0_.title as title7_0_0_ from job job0_ where job0_.id=5",
		"2024-04-08T14:10:41.836540Z	 2639 Query	select job0_.id as id1_0_0_, job0_.description as descript2_0_0_, job0_.publish_at as publish_3_0_0_, job0_.publish_trials as publish_4_0_0_, job0_.published_timestamp as publishe5_0_0_, job0_.tags as tags6_0_0_, job0_.title as title7_0_0_ from job job0_ where job0_.id=5",
		"2024-04-08T14:10:41.837518Z	 2639 Query	update job set description='World, X', publish_at='2024-04-08 16:10:21', publish_trials=0, published_timestamp=null, tags='', title='Hello' where id=5",
		"2024-04-08T14:10:41.839523Z	 2639 Query	select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_",
		"2024-04-08T14:10:42.547042Z	 2639 Query	update job set description='World, X', publish_at='2024-04-08 16:10:21', publish_trials=1, published_timestamp='2024-04-08 16:10:42.545649', tags='', title='Hello' where id=5",
	}
	verifications := []string{
		"2024-04-08T09:36:13.070009Z	 2549 Query	select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ order by job0_.publish_at desc",
		"2024-04-08T12:48:25.907355Z	 2609 Query	insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World', '2024-04-08 14:48:15', 0, null, '', 'Hello', 2)",
		"2024-04-08T12:48:25.908907Z	 2609 Query	select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ order by job0_.publish_at desc",
		"2024-04-08T12:48:26.838373Z	 2609 Query	update job set description='World', publish_at='2024-04-08 14:48:15', publish_trials=1, published_timestamp='2024-04-08 14:48:26.836784', tags='', title='Hello' where id=2",
		"2024-04-08T14:12:59.196445Z	 2639 Query	select job0_.id as id1_0_0_, job0_.description as descript2_0_0_, job0_.publish_at as publish_3_0_0_, job0_.publish_trials as publish_4_0_0_, job0_.published_timestamp as publishe5_0_0_, job0_.tags as tags6_0_0_, job0_.title as title7_0_0_ from job job0_ where job0_.id=6",
		"2024-04-08T14:13:04.090246Z	 2639 Query	select job0_.id as id1_0_0_, job0_.description as descript2_0_0_, job0_.publish_at as publish_3_0_0_, job0_.publish_trials as publish_4_0_0_, job0_.published_timestamp as publishe5_0_0_, job0_.tags as tags6_0_0_, job0_.title as title7_0_0_ from job job0_ where job0_.id=6",
		"2024-04-08T14:13:04.091495Z	 2639 Query	update job set description='World, X', publish_at='2024-04-08 16:12:46', publish_trials=0, published_timestamp=null, tags='', title='Hello' where id=6",
		"2024-04-08T14:13:04.094251Z	 2639 Query	select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_",
		"2024-04-08T14:13:04.667827Z	 2639 Query	update job set description='World, X', publish_at='2024-04-08 16:12:46', publish_trials=1, published_timestamp='2024-04-08 16:13:04.666452', tags='', title='Hello' where id=6",
	}
	expectationSource := adapter.NewMemExpectationSource(expectations)
	verificationSource := adapter.NewMemExpectationSource(verifications)

	logs := []string{
		"2024-04-08T09:33:15.070009Z	 2549 Query	select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ order by job0_.publish_at desc",
		"2024-04-08T12:50:59.605638Z	 2609 Query	insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World', '2024-04-08 14:50:20', 0, null, '', 'Hello', 3)",
		"2024-04-08T12:50:59.607117Z	 2609 Query	select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ order by job0_.publish_at desc",
		"2024-04-08T12:51:00.619360Z	 2609 Query	update job set description='World', publish_at='2024-04-08 14:50:20', publish_trials=1, published_timestamp='2024-04-08 14:51:00.618489', tags='', title='Hello' where id=3",
		"2024-04-08T14:14:23.063252Z	 2639 Query	select job0_.id as id1_0_0_, job0_.description as descript2_0_0_, job0_.publish_at as publish_3_0_0_, job0_.publish_trials as publish_4_0_0_, job0_.published_timestamp as publishe5_0_0_, job0_.tags as tags6_0_0_, job0_.title as title7_0_0_ from job job0_ where job0_.id=8",
		"2024-04-08T14:14:31.973778Z	 2639 Query	select job0_.id as id1_0_0_, job0_.description as descript2_0_0_, job0_.publish_at as publish_3_0_0_, job0_.publish_trials as publish_4_0_0_, job0_.published_timestamp as publishe5_0_0_, job0_.tags as tags6_0_0_, job0_.title as title7_0_0_ from job job0_ where job0_.id=8",
		"2024-04-08T14:14:31.974878Z	 2639 Query	update job set description='World, X', publish_at='2024-04-08 16:14:15', publish_trials=0, published_timestamp=null, tags='', title='Hello' where id=8",
		"2024-04-08T14:14:31.976795Z	 2639 Query	select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_",
		"2024-04-08T14:14:32.571579Z	 2639 Query	update job set description='World, X', publish_at='2024-04-08 16:14:15', publish_trials=1, published_timestamp='2024-04-08 16:14:32.570201', tags='', title='Hello' where id=8",

		"STOP",
	}

	databaseLog := adapter.NewMemSQLLog(logs)
	timer := adapter.MockTimer{}
	listener = NewListener(c, timer, databaseLog, expectationSource, verificationSource)
	listener.Start()
	listener.Stop()
	results := listener.GetResults()
	if len(results) > 0 {
		for _, v := range results {
			log.Println(v)

		}
	}
	assert.Len(t, results, 0)
}
