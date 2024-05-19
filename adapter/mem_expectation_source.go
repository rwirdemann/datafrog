package adapter

import (
	"github.com/rwirdemann/datafrog/internal/datafrog"
)

type MemExpectationSource struct {
	testcase datafrog.Testcase
}

func NewMemExpectationSource(expectations []datafrog.Expectation) *MemExpectationSource {
	return &MemExpectationSource{testcase: datafrog.Testcase{Expectations: expectations}}
}

func (es MemExpectationSource) Get() datafrog.Testcase {
	return es.testcase
}

func (es *MemExpectationSource) Write(testcase datafrog.Testcase) error {
	es.testcase = testcase
	return nil
}
