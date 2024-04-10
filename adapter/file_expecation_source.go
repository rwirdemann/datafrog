package adapter

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/rwirdemann/databasedragon/matcher"
)

type FileExpectationSource struct {
	expectations []string
}

func NewFileExpectationSource(filename string) *FileExpectationSource {
	expectations, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	split := strings.Split(string(expectations), "\n")
	var all []string
	for _, v := range split {
		if strings.Trim(v, " ") != "" {
			all = append(all, v)
		}
	}
	return &FileExpectationSource{expectations: all}
}

func (s *FileExpectationSource) GetAll() []string {
	return s.expectations
}

// RemoveFirst removes the first expectation if it matches the pattern.
// Returns an error if no remaining expectations are left or if the first
// expectations doesn't match the pattern.
func (s *FileExpectationSource) RemoveFirst(pattern string) error {
	if len(s.expectations) == 0 {
		return errors.New("list of expectations is empty")
	}

	if !matcher.NewPattern(pattern).MatchesAllConditions(s.expectations[0]) {
		return errors.New("first expectation didn't match expected pattern")
	}

	s.expectations = s.expectations[1:]
	return nil
}
