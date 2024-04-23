package adapter

import "github.com/rwirdemann/databasedragon/matcher"

type MemExpectationSource struct {
	expecations []matcher.Expectation
}

func NewMemExpectationSource(expectations []matcher.Expectation) *MemExpectationSource {
	return &MemExpectationSource{expecations: expectations}
}

func (es *MemExpectationSource) GetAll() []matcher.Expectation {
	return es.expecations
}

func (es *MemExpectationSource) WriteAll() {
}
