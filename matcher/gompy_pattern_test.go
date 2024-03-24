package matcher

import "testing"

func TestGompyPattern_MatchesAllConditions(t *testing.T) {
	tests := []struct {
		name    string
		actual  string
		pattern string
		want    bool
	}{
		{
			name:    "matches",
			actual:  "insert into job (description, publish_at, publish_trials, published_timestamp, tags, title) values ('Hello', '2024-03-24 11:46:46', 0, null, '', 'World')",
			pattern: "insert into job (description, publish_at, publish_trials, published_timestamp, tags, title) values ('Hello', '<DATE_STR>', 0, null, '', 'World')",
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := GompyPattern{
				Include: tt.pattern,
			}
			if got := p.MatchesAllConditions(tt.actual); got != tt.want {
				t.Errorf("MatchesAllConditions() = %v, want %v", got, tt.want)
			}
		})
	}
}
