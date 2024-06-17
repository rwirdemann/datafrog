package mocks

import (
	"errors"
	"github.com/rwirdemann/datafrog/pkg/df"
)

type TestRepository struct {
	Testcases []df.Testcase
}

func (r TestRepository) All() ([]df.Testcase, error) {
	return r.Testcases, nil
}

func (r TestRepository) Get(filename string) (df.Testcase, error) {
	for _, tc := range r.Testcases {
		if tc.Name == filename {
			return tc, nil
		}
	}
	return df.Testcase{}, errors.New("testcase not found")
}

func (r TestRepository) Exists(filename string) bool {
	_, err := r.Get(filename)
	if err != nil {
		return false
	}
	return true
}
