package mocks

import (
	"github.com/rwirdemann/datafrog/pkg/df"
)

type ExpectationSource struct {
	testcase df.Testcase
}

func NewExpectationSource(expectations []df.Expectation) *ExpectationSource {
	return &ExpectationSource{testcase: df.Testcase{Expectations: expectations}}
}

func (es ExpectationSource) Get() df.Testcase {
	return es.testcase
}

func (es *ExpectationSource) Write(testcase df.Testcase) error {
	es.testcase = testcase
	return nil
}
