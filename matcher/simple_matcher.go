package matcher

import (
	"github.com/rwirdemann/texttools/config"
	"log"
)

type SimpleMatcher struct {
	config config.Config
}

func (m SimpleMatcher) MatchingPattern(s string) string {
	for _, pattern := range m.config.Patterns {
		if NewPattern(pattern).MatchesAllConditions(s) {
			return pattern
		}
	}
	log.Fatalf("Matching pattern not found in '%s'", s)
	return ""
}

func NewPatternMatcher(config config.Config) Matcher {
	return SimpleMatcher{config: config}
}

func (m SimpleMatcher) MatchesAny(s string) bool {
	for _, p := range m.config.Patterns {
		if NewPattern(p).MatchesAllConditions(s) {
			return true
		}
	}
	return false
}
