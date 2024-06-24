package mocks

import (
	"errors"
	"github.com/rwirdemann/datafrog/pkg/df"
)

type TestRepository struct {
	Testcases []df.Testcase
}

func (r *TestRepository) Delete(testname string) error {
	panic("implement me")
}

func (r *TestRepository) Write(_ string, testcase df.Testcase) error {
	r.Testcases = append(r.Testcases, testcase)
	return nil
}

func (r *TestRepository) All() ([]df.Testcase, error) {
	return r.Testcases, nil
}

func (r *TestRepository) Get(filename string) (df.Testcase, error) {
	for _, tc := range r.Testcases {
		if tc.Name == filename {
			return tc, nil
		}
	}
	return df.Testcase{}, errors.New("testcase not found")
}

func (r *TestRepository) Exists(filename string) bool {
	_, err := r.Get(filename)
	if err != nil {
		return false
	}
	return true
}
