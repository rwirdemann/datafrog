package adapter

import "github.com/rwirdemann/databasedragon/matcher"

type MemExpectationSource struct {
	expecations []matcher.Expectation
}

func NewMemExpectationSource(expecations []matcher.Expectation) *MemExpectationSource {
	return &MemExpectationSource{expecations: expecations}
}

func (es *MemExpectationSource) GetAll() []matcher.Expectation {
	return es.expecations
}

func (es *MemExpectationSource) WriteAll(expectations []matcher.Expectation) {
}
