package mysql

import (
	"github.com/rwirdemann/datafrog/internal/datafrog"
	"strings"
)

type Tokenizer struct {
}

// Tokenize cuts timestamp and any additional characters that not are part of
// the plain sql statement from s. The cleaned statements is split by spaces
// into single tokens afterward.
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
