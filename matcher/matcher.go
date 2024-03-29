package matcher

type Matcher interface {
	MatchesPattern(s string) bool
	MatchesExactly(s1 string, s2 string) bool
}
