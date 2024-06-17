package df

import (
	"strings"
)

type Pattern struct {
	Include string
	Exclude string
}

func NewPattern(s string) Pattern {
	e := strings.Split(s, "!")
	if len(e) == 2 {
		return Pattern{
			Include: e[0],
			Exclude: e[1],
		}
	}
	return Pattern{
		Include: s,
		Exclude: "",
	}
}

func MatchesPattern(patterns []string, s string) (bool, string) {
	for _, p := range patterns {
		if NewPattern(p).matches(s) {
			return true, p
		}
	}
	return false, ""
}

func (p Pattern) matchesInclude(s string) bool {
	return strings.Contains(strings.ToUpper(s), strings.ToUpper(p.Include))
}

func (p Pattern) matchesExclude(s string) bool {
	return len(p.Exclude) > 0 && strings.Contains(strings.ToUpper(s), strings.ToUpper(p.Exclude))
}

func (p Pattern) matches(s string) bool {
	return p.matchesInclude(s) && !p.matchesExclude(s)
}
