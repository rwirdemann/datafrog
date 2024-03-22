package matcher

import (
	"github.com/rwirdemann/texttools/config"
	"strings"
)

type SimpleMatcher struct {
	config config.Config
}

func (m SimpleMatcher) MatchingPattern(s string) string {
	for _, pattern := range m.config.Patterns {
		if strings.Contains(s, pattern) {
			return pattern
		}
	}
	return ""
}

func NewPatternMatcher(config config.Config) Matcher {
	return SimpleMatcher{config: config}
}

func (m SimpleMatcher) MatchesAny(s string) bool {
	for _, p := range m.config.Patterns {
		p := NewPattern(p)
		if p.MatchesInclude(s) && !p.MatchesExclude(s) {
			return true
		}
	}
	return false
}
