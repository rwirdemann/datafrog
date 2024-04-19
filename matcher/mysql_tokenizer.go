package matcher

type MySQLTokenizer struct {
}

// Tokenize cuts timestamp and any additional characters that not are part of
// the plain sql statement from s. The cleaned statements is split by spaces
// into single tokens afterward.
func (m MySQLTokenizer) Tokenize(s string, patterns []string) []string {
	return Tokenize(normalize(cutPrefix(s, patterns), patterns))
}
