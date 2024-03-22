package matcher

import "strings"

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

func (p Pattern) MatchesInclude(s string) bool {
	return strings.Contains(s, p.Include)
}

func (p Pattern) MatchesExclude(s string) bool {
	return len(p.Exclude) > 0 && strings.Contains(s, p.Exclude)
}
