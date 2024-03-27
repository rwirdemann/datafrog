package matcher

import (
	"github.com/rwirdemann/databasedragon/config"
	"testing"
)

func TestMatchesAny(t *testing.T) {
	c := config.NewConfig("simple_config.json")

	m := NewPatternMatcher(c)
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"matches insert", "2024-03-21T11:28:39.975400Z\t  241 Query\tinsert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('fdfdf', '2024-03-21 12:28:36', 0, null, '', 'fdff', 16)\n", true},
		{"not matches insert", "2024-03-21T11:28:39.975400Z\t  241 Query\tinsert into user (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('fdfdf', '2024-03-21 12:28:36', 0, null, '', 'fdff', 16)\n", false},
		{"matches update", "2024-03-21T11:28:40.433827Z\t  241 Query\tupdate job set description='fdfdf', publish_at='2024-03-21 12:28:36', publish_trials=1, published_timestamp='2024-03-21 12:28:40.432771', tags='', title='fdff' where id=16\n", true},
		{"not matche update", "2024-03-21T11:28:40.433827Z\t  241 Query\tupdates job set description='fdfdf', publish_at='2024-03-21 12:28:36', publish_trials=1, published_timestamp='2024-03-21 12:28:40.432771', tags='', title='fdff' where id=16\n", false},
		{"matches select", "2024-03-22T11:05:39.164989Z\t  311 Query\tselect job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ order by job0_.publish_at desc\n", true},
		{"exclude", "2024-03-22T11:07:13.295519Z\t  321 Query\tselect job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ where (job0_.published_timestamp is null) and job0_.publish_at<'2024-03-22 12:07:13.291579' and job0_.publish_trials<1\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.MatchesAny(tt.s); got != tt.want {
				t.Errorf("MatchesAny() = %v, want %v", got, tt.want)
			}
		})
	}
}
