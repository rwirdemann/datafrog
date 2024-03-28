package matcher

import (
	"log"

	"github.com/rwirdemann/databasedragon/config"
)

type GompyMatcher struct {
	config config.Config
}

func NewGompyMatcher(config config.Config) Matcher {
	return GompyMatcher{config: config}
}

func (m GompyMatcher) MatchingPattern(s string) string {
	for _, pattern := range m.config.Patterns {
		if NewPattern(pattern).MatchesAllConditions(s) {
			return pattern
		}
	}
	log.Fatalf("Matching pattern not found in '%s'", s)
	return ""
}

func (m GompyMatcher) MatchesAny(s string) bool {
	for _, p := range m.config.Patterns {
		if NewPattern(p).MatchesAllConditions(s) {
			return true
		}
	}
	return false
}

func (p GompyMatcher) MatchesExactly(s1 string, s2 string) bool {
	return false
}
