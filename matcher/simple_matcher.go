package matcher

import (
	"github.com/rwirdemann/texttools/config"
	"strings"
)

type SimpleMatcher struct {
	config config.Config
}

func NewSimpleMatcher(config config.Config) Matcher {
	return SimpleMatcher{config: config}
}

func (m SimpleMatcher) MatchesAny(s string) bool {
	for _, pattern := range m.config.Include {
		if strings.Contains(s, pattern) {
			return true
		}
	}
	return false
}
