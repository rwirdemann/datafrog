package mocks

import (
	"github.com/rwirdemann/datafrog/pkg/df"
)

type TestRepository struct {
	Testcases []df.Testcase
}

func (r TestRepository) All() ([]df.Testcase, error) {
	return r.Testcases, nil
}
