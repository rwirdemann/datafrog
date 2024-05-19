package postgres

import (
	"github.com/rwirdemann/datafrog/internal/datafrog"
	"strings"
)

// Tokenizer tokenizes PostgreSQL log entries. PostgreSQL configuration
// settings:
//
//	log_destination = 'stderr,csvlog'
//
// PostgreSQL splits single sql statements into two parts:
//
//	1: select * from job where job0_.publish_at<$1 and job0_.publish_trials<$2
//	2: parameters: $1 = '2024-04-19 10:07:38.543981', $2 = '1'
type Tokenizer struct {
}

func (m Tokenizer) Tokenize(s string, patterns []string) []string {
	return datafrog.Tokenize(normalize(cutPrefix(s, patterns), patterns))
}

func normalize(s string, patterns []string) string {
	result := cutPrefix(s, patterns)
	result = strings.TrimSuffix(result, "\n")
	return result
}

func cutPrefix(s string, patterns []string) string {
	for _, p := range patterns {
		idx := strings.Index(s, datafrog.NewPattern(p).Include)
		if idx > -1 {
			return s[idx:]
		}
	}
	return s
}
