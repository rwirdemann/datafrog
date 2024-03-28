package matcher

import (
	"log"
	"regexp"
	"strings"

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
	s1 = cutPrefix(s1, m.config.Patterns)
	s1 = strings.TrimSuffix(s1, "\n")
	s2 = cutPrefix(s2, m.config.Patterns)
	s2 = strings.TrimSuffix(s2, "\n")

	r := regexp.MustCompile(`([A-Za-z0-9]+(-[A-Za-z0-9]+)+) ([A-Za-z0-9]+(:[A-Za-z0-9]+)+)(\.[0-9]+)?`)
	s1 = r.ReplaceAllString(s1, "<DATE_STR>")
	s2 = r.ReplaceAllString(s2, "<DATE_STR>")

	r2 := regexp.MustCompile(`(, [0-9]+)`)
	s1 = r2.ReplaceAllString(s1, "<ID>")
	s2 = r2.ReplaceAllString(s2, "<ID>")

	r3 := regexp.MustCompile(`(=[0-9]+)`)
	s1 = r3.ReplaceAllString(s1, "<ID>")
	s2 = r3.ReplaceAllString(s2, "<ID>")

	return s1 == s2
}

func useRegex(s string) bool {
	re := regexp.MustCompile("([A-Za-z0-9]+(-[A-Za-z0-9]+)+) ([A-Za-z0-9]+(:[A-Za-z0-9]+)+)")
	return re.MatchString(s)
}

func cutPrefix(s string, patterns []string) string {
	for _, p := range patterns {
		idx := strings.Index(s, NewPattern(p).Include)
		if idx > -1 {
			return s[idx:]
		}
	}
	return s
}
