package matcher

type Matcher interface {
	MatchesAny(s string) bool
	MatchingPattern(s string) string
}
