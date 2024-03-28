package matcher

import (
	"fmt"
	"testing"

	"github.com/rwirdemann/databasedragon/config"
)

func TestMatchesExactly(t *testing.T) {
	c := config.NewConfig("simple_config.json")

	m := NewDynamicDataMatcher(c)
	tests := []struct {
		s1   string
		s2   string
		want bool
	}{
		{
			"insert into job (description, publish_at) values ('Hello', '2024-03-24 11:46:46')",
			"insert into job (description, publish_at) values ('Hello', '2024-03-24 12:48:22')",
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
