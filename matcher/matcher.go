package matcher

type Matcher interface {
	MatchesAny(s string) bool
	MatchingPattern(s string) string
	MatchesExactly(s1 string, s2 string) bool
}
