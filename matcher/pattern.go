package matcher

import (
	"github.com/rwirdemann/datafrog/internal/datafrog"
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

func (p Pattern) MatchesInclude(s string) bool {
	return strings.Contains(s, p.Include)
}

func (p Pattern) MatchesExclude(s string) bool {
	return len(p.Exclude) > 0 && strings.Contains(s, p.Exclude)
}

func (p Pattern) MatchesAllConditions(s string) bool {
	return p.MatchesInclude(s) && !p.MatchesExclude(s)
}

func MatchesPattern(c datafrog.Config, s string) (bool, string) {
	for _, p := range c.Patterns {
		if NewPattern(p).MatchesAllConditions(s) {
			return true, p
		}
	}
	return false, ""
}
