package matcher

import "testing"

func TestGompyPattern_MatchesAllConditions(t *testing.T) {
	pattern := "insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('Hello', '<DATE_STR>', 0, null, '', 'World', <ID>)"
	tests := []struct {
		name    string
		actual  string
		pattern string
		want    bool
	}{
		{
			name:    "m3",
			actual:  "2024-03-24T10:46:52.226470Z 821 Query insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('Hello', '2024-03-24 11:46:46', 0, null, '', 'World', 14)",
			pattern: pattern,
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
