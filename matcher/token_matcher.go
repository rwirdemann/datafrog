package matcher

import (
	"log"
	"strings"

	"github.com/rwirdemann/databasedragon/config"
)

type TokenMatcher struct {
	config       config.Config
	expectations []Expectation
}

func NewTokenMatcher(c config.Config, expecations, verifications []string) TokenMatcher {
	tm := TokenMatcher{config: c}
	for i, v := range expecations {
		if strings.Trim(v, " ") != "" {
			e := NewExpectation(normalize(v, c.Patterns), normalize(verifications[i], c.Patterns))
			tm.expectations = append(tm.expectations, e)
		}
	}
	return tm
}

func (t *TokenMatcher) Matches(actual string) int {
	normalized := normalize(actual, t.config.Patterns)
	for i, v := range t.expectations {
		if v.Equal(normalized) {
			return i
		}
	}
	return -1
}

func (t *TokenMatcher) RemoveExpectation(i int) {
	t.expectations = append(t.expectations[:i], t.expectations[i+1:]...)
}

func (t *TokenMatcher) PrintResults() {
	if len(t.expectations) == 0 {
		log.Printf("All expectations met!")
	} else {
		log.Printf("Failed due to missing expectations! Missing: %d", len(t.expectations))
	}
}

func (t *TokenMatcher) GetResults() []Expectation {
	return t.expectations
}

func normalize(s string, patterns []string) string {
	result := cutPrefix(s, patterns)
	result = strings.TrimSuffix(result, "\n")
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
