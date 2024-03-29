package matcher

import (
	"log"
	"regexp"
	"strings"

	"github.com/agnivade/levenshtein"
	"github.com/rwirdemann/databasedragon/config"
)

// Min ratio to consider two strings as equal according to levenshtein string comparison.
const minLevenshteinTresholdRatio = 0.99

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

func (m LevenshteinMatcher) MatchesPattern(s string) bool {
	for _, p := range m.config.Patterns {
		if NewPattern(p).MatchesAllConditions(s) {
			return true
		}
	}
	return false
}

func (m LevenshteinMatcher) MatchesExactly(recorded string, expecation string) bool {

	// Quickcheck: Does the recorded line contain one of the recording patterns?
	if !m.MatchesPattern(recorded) {
		return false
	}

	recorded = cutPrefix(recorded, m.config.Patterns)
	recorded = strings.TrimSuffix(recorded, "\n")
	expecation = cutPrefix(expecation, m.config.Patterns)
	expecation = strings.TrimSuffix(expecation, "\n")

	r := regexp.MustCompile(`([A-Za-z0-9]+(-[A-Za-z0-9]+)+) ([A-Za-z0-9]+(:[A-Za-z0-9]+)+)(\.[0-9]+)?`)
	recorded = r.ReplaceAllString(recorded, "<DATE_STR>")
	expecation = r.ReplaceAllString(expecation, "<DATE_STR>")

	distance := float64(levenshtein.ComputeDistance(recorded, expecation))
	sum := float64(len(recorded)) + float64(len(expecation))

	ratio := (sum - distance) / sum

	log.Printf("EXPECTATION: %s\n", expecation)
	log.Printf("RECORDED   : %s\n", recorded)

	match := ratio > minLevenshteinTresholdRatio
	if match {
		log.Printf("LEVEN RATIO: %f => Expectation met", ratio)
	} else {
		log.Printf("LEVEN RATIO: %f => Expectation failed", ratio)
	}
	log.Println("---------------------------------------------------------------------------------")
	return match
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
