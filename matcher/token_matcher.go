package matcher

import (
	"github.com/rwirdemann/databasedragon/config"
)

type TokenMatcher struct {
	config       config.Config
	expectations []Expectation
}

func NewTokenMatcher(c config.Config, expecations, verifications []string) TokenMatcher {
	tm := TokenMatcher{config: c}
	for i, v := range expecations {
		e := NewExpectation(normalize(v, c.Patterns), normalize(verifications[i], c.Patterns))
		tm.expectations = append(tm.expectations, e)
	}
	return tm
}

func (t TokenMatcher) Matches(actual string) bool {
	normalized := normalize(actual, t.config.Patterns)
	for _, v := range t.expectations {
		if v.Equal(normalized) {
			return true
		}
	}
	return false
}
