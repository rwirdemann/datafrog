package matcher

import (
	"regexp"
	"strings"
)

type GompyPattern struct {
	Include string
	Exclude string
}

func NewGompyPattern(s string) GompyPattern {
	e := strings.Split(s, "!")
	if len(e) == 2 {
		return GompyPattern{
			Include: e[0],
			Exclude: e[1],
		}
	}
	return GompyPattern{
		Include: s,
		Exclude: "",
	}
}

func (p GompyPattern) MatchesInclude(s string) bool {
	return strings.Contains(s, p.Include)
}

func (p GompyPattern) MatchesExclude(s string) bool {
	return len(p.Exclude) > 0 && strings.Contains(s, p.Exclude)
}

func (p GompyPattern) MatchesAllConditions(s string) bool {
	regex := strings.ReplaceAll(p.Include, "(", "\\(")
	regex = strings.ReplaceAll(regex, ")", "\\)")
	regex = strings.ReplaceAll(regex, "<DATE_STR>", "([0-9\\-]+ [0-9\\:]+)")
	r, _ := regexp.Compile(regex)
	return r.MatchString(s)
}
