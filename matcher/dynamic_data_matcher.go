package matcher

import (
	"log"
	"regexp"

	"github.com/rwirdemann/databasedragon/config"
)

type DynamicDataMatcher struct {
	config config.Config
}

func NewDynamicDataMatcher(config config.Config) DynamicDataMatcher {
	return DynamicDataMatcher{config: config}
}

func (m DynamicDataMatcher) MatchingPattern(s string) string {
	for _, pattern := range m.config.Patterns {
		if NewPattern(pattern).MatchesAllConditions(s) {
			return pattern
		}
	}
	log.Fatalf("Matching pattern not found in '%s'", s)
	return ""
}

func (m DynamicDataMatcher) MatchesAny(s string) bool {
	for _, p := range m.config.Patterns {
		if NewPattern(p).MatchesAllConditions(s) {
			return true
		}
	}
	return false
}

func (m DynamicDataMatcher) MatchesExactly(s1 string, s2 string) bool {
	r := regexp.MustCompile(`([0-9\\-]+ [0-9\\:]+)`)
	s1 = r.ReplaceAllString(s1, "<DATE_STR>")
	s2 = r.ReplaceAllString(s2, "<DATE_STR>")
	return s1 == s2
}
