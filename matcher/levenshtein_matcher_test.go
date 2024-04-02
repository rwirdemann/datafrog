package matcher

import (
	"fmt"
	"testing"

	"github.com/rwirdemann/databasedragon/config"
	"github.com/stretchr/testify/assert"
)

func TestMatchesExactly(t *testing.T) {
	c := config.NewConfig("config.json")

	m := NewLevenshteinMatcher(c)
	tests := []struct {
		s1   string
		s2   string
		want bool
	}{
		{
			"update job set description='World', publish_at='2024-03-28 10:57:27', publish_trials=1, published_timestamp='2024-03-28 10:57:36.261562', tags='', title='Hello' where id=11",
			"update job set description='World', publish_at='2024-03-28 10:58:27', publish_trials=1, published_timestamp='2024-03-28 10:58:36.261563', tags='', title='Hello' where id=12",
			true,
		},
		{
			"insert into job (description, publish_at) values ('Hello', '2024-03-24 11:46:46')",
			"insert into job (description, publish_at) values ('Hello', '2024-03-24 12:48:22')",
			true,
		},
		{
			"2024-03-24T10:46:52.226470Z 821 Query insert into job (description, publish_at) values ('Hey', '2024-03-24 11:46:46')",
			"2024-03-24T11:46:52.226470Z 821 Query insert into job (description, publish_at) values ('Hey', '2024-03-24 12:48:22')",
			true,
		},
		{
			"2024-03-28T08:50:25.126344Z	  599 Query	select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ order by job0_.publish_at desc",
			"2024-03-28T08:53:25.126344Z	  600 Query	select job0_.id as id1_0_, job0_.description as descript2_0_, job0_.publish_at as publish_3_0_, job0_.publish_trials as publish_4_0_, job0_.published_timestamp as publishe5_0_, job0_.tags as tags6_0_, job0_.title as title7_0_ from job job0_ order by job0_.publish_at desc",
			true,
		},
		{
			"insert into job (description, id) values ('World', 11)",
			"insert into job (description, id) values ('World', 12)",
			true,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			if got := m.MatchesExactly(tt.s1, tt.s2); got != tt.want {
				t.Errorf("MatchesExactly() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalize(t *testing.T) {
	tests := []struct {
		actual   string
		expected string
	}{
		{
			actual:   "2024-03-29T11:15:00.509261Z	  879 Query	select job0_.id",
			expected: "select job0_.id",
		},
		{
			actual:   "2024-03-29T11:15:00.509261Z	  879 Query	select job0_.id\n",
			expected: "select job0_.id",
		},
		{
			actual:   "2024-03-29T11:15:08.238360Z	  879 Query	update job set description='w', publish_at='2024-03-29 12:15:01'",
			expected: "update job set description='w', publish_at='<DATE_STR>'",
		},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.expected, normalize(tc.actual, []string{"select", "update"}))
	}
}
