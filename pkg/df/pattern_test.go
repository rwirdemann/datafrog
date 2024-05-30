package df

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var config = Config{Patterns: []string{
	"update job",
	"insert into",
	"delete",
	"select job!publish_trials<1",
}}

func TestMatchesPattern(t *testing.T) {
	tests := []struct {
		s               string
		expectedMatches bool
		expectedPattern string
		name            string
	}{
		{name: "matches insert", expectedPattern: "insert into", expectedMatches: true, s: "insert into"},
		{name: "considers exclude", expectedPattern: "", expectedMatches: false, s: "select job and publish_trials<1"},
		{name: "ignore case", expectedPattern: "insert into", expectedMatches: true, s: "INSERT INTO JOB"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			matches, pattern := MatchesPattern(config, test.s)
			assert.Equal(t, test.expectedMatches, matches)
			assert.Equal(t, test.expectedPattern, pattern)
		})
	}
}
