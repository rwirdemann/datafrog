package matcher

import (
	"github.com/rwirdemann/texttools/config"
	"testing"
)

func TestMatchesAny(t *testing.T) {
	c := config.NewConfig("test_config.json")

	m := NewSimpleMatcher(c)
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"matches insert", "2024-03-21T11:28:39.975400Z\t  241 Query\tinsert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('fdfdf', '2024-03-21 12:28:36', 0, null, '', 'fdff', 16)\n", true},
		{"not matches insert", "2024-03-21T11:28:39.975400Z\t  241 Query\tinsert into user (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('fdfdf', '2024-03-21 12:28:36', 0, null, '', 'fdff', 16)\n", false},
		{"matches update", "2024-03-21T11:28:40.433827Z\t  241 Query\tupdate job set description='fdfdf', publish_at='2024-03-21 12:28:36', publish_trials=1, published_timestamp='2024-03-21 12:28:40.432771', tags='', title='fdff' where id=16\n", true},
		{"not matche update", "2024-03-21T11:28:40.433827Z\t  241 Query\tupdates job set description='fdfdf', publish_at='2024-03-21 12:28:36', publish_trials=1, published_timestamp='2024-03-21 12:28:40.432771', tags='', title='fdff' where id=16\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.MatchesAny(tt.s); got != tt.want {
				t.Errorf("MatchesAny() = %v, want %v", got, tt.want)
			}
		})
	}
}
