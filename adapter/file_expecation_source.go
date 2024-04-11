package adapter

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/rwirdemann/databasedragon/matcher"
)

// A FileExpectationSource reads raw (plain SQL or other log entries)
// expectations from a file. It is used by verification runs to verfify that the
// same log entries are written and to build its differences.
//
// The file based expectations are stored in expectations. The idea is to remove
// always the first item from this list when the verfication process has
// retrieved a pattern matching line from the log. The matching pattern is
// applied to the first entry in expectations. The entry is removed if it
// matches and the verfication is successfull, after all expectations have been
// removed. See RemoveFirst for details.
type FileExpectationSource struct {
	expectations []string
}

// NewFileExpectationSource creates a new NewFileExpectationSource that reads
// its expectations from filename.
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

// GetAll returns all expectations.
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
