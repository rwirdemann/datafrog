package matcher

import "strings"

type MySQLTokenizer struct {
}

// Tokenize cuts timestamp and any additional characters that not are part of
// the plain sql statement from s. The cleaned statements is split by spaces
// into single tokens afterward.
func (m MySQLTokenizer) Tokenize(s string, patterns []string) []string {
	return Tokenize(normalize(cutPrefix(s, patterns), patterns))
}

func normalize(s string, patterns []string) string {
	result := cutPrefix(s, patterns)
	result = strings.TrimSuffix(result, "\n")
	return result
}

func cutPrefix(s string, patterns []string) string {
	for _, p := range patterns {
		idx := strings.Index(s, NewPattern(p).Include)
		if idx > -1 {
			return s[idx:]
		}
	}
	return s
}
