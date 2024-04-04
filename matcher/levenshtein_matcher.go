package matcher

import (
	"log"
	"regexp"
	"strings"

	"github.com/agnivade/levenshtein"
	"github.com/rwirdemann/databasedragon/config"
)

// The LevenshteinMatcher fulfills the Matcher interface and uses the Levenshtein algorithms for
// comparing SQL statements.
type LevenshteinMatcher struct {
	config config.Config
}

func NewLevenshteinMatcher(config config.Config) LevenshteinMatcher {
	return LevenshteinMatcher{config: config}
}

func (m LevenshteinMatcher) MatchingPattern(s string) string {
	for _, pattern := range m.config.Patterns {
		if NewPattern(pattern).MatchesAllConditions(s) {
			return pattern
		}
	}
	log.Fatalf("Matching pattern not found in '%s'", s)
	return ""
}

func (m LevenshteinMatcher) MatchesPattern(s string) (bool, Pattern) {
	for _, p := range m.config.Patterns {
		if NewPattern(p).MatchesAllConditions(s) {
			return true, NewPattern(p)
		}
	}
	return false, Pattern{}
}

func (m LevenshteinMatcher) MatchesExactly(recorded string, expecation string) bool {

	// Remove dynamic and other noisy data
	recorded = normalize(recorded, m.config.Patterns)
	expecation = normalize(expecation, m.config.Patterns)

	distance := levenshtein.ComputeDistance(recorded, expecation)

	var Green = "\033[32m"
	var Red = "\033[31m"
	match := distance <= m.config.MaxLevenshteinDistance
	if match {
		log.Printf(Green+"Levenshtein Distance: %d => Expectation met", distance)
	} else {
		log.Printf(Red+"Levenshtein Distance: %d => Expectation failed", distance)
	}
	return match
}

func normalize(s string, patterns []string) string {
	result := cutPrefix(s, patterns)
	result = strings.TrimSuffix(result, "\n")
	timeRegex := regexp.MustCompile(`([A-Za-z0-9]+(-[A-Za-z0-9]+)+) ([A-Za-z0-9]+(:[A-Za-z0-9]+)+)(\.[0-9]+)?`)
	result = timeRegex.ReplaceAllString(result, "<DATE_STR>")

	return result
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
