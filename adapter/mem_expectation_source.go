package adapter

import (
	"github.com/rwirdemann/databasedragon/app/domain"
)

type MemExpectationSource struct {
	expecations []domain.Expectation
}

func NewMemExpectationSource(expectations []domain.Expectation) *MemExpectationSource {
	return &MemExpectationSource{expecations: expectations}
}

func (es *MemExpectationSource) GetAll() []domain.Expectation {
	return es.expecations
}

func (es *MemExpectationSource) WriteAll() {
}
