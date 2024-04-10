package adapter

import (
	"errors"

	"github.com/rwirdemann/databasedragon/matcher"
)

type MemExpectationSource struct {
	expecations []string
}

func NewMemExpectationSource(expecations []string) *MemExpectationSource {
	return &MemExpectationSource{expecations: expecations}
}

func (es *MemExpectationSource) GetAll() []string {
	return es.expecations
}

// RemoveFirst removes the first expectation if it matches the pattern.
// Returns an error if no remaining expectations are left or if the first
// expectations doesn't match the pattern.
func (es *MemExpectationSource) RemoveFirst(pattern string) error {
	if len(es.expecations) == 0 {
		return errors.New("list of expectations is empty")
	}

	if !matcher.NewPattern(pattern).MatchesAllConditions(es.expecations[0]) {
		return errors.New("first expectation didn't match expected pattern")
	}

	es.expecations = es.expecations[1:]
	return nil
}
