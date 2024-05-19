package datafrog

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Expectation represents a log entry that is expected to reappear within
// succeeding test runs. The log entry must be presented as an array of tokens.
// Example: "2024-04-16 13:53:53.277 select * from job where id=12" should
// become ["select", "*", "from", "job", "where", "id=12"] within an Expectation
// instance. An expectation matches a Pattern. All leading characters before
// this pattern must be removed from the raw entry before building the token set.
//
// An expectation that could be verified (reappeared) with in a test run becomes
// "Fulfilled" and gets an by one increased "Verified" count. Expectations
// should be persisted between test runs. Thus their Verified counter increases
// over time and the overall test quality gains.
type Expectation struct {
	Uuid      string   `json:"uuid"`
	Tokens    []string `json:"tokens"`
	Pattern   string
	Fulfilled bool
	Verified  int

	IgnoreDiffs []int `json:"ignoreDiffs"` // indizes of tokens allowed to deviate when comparing two Expectations
}

// Equal compares e's tokens with the given tokens. The tokens sets are equal if
// they have the same length, all their elements are equal or if two elements
// are unequal but their index is contained in IgnoreDiffs.
func (e Expectation) Equal(tokens []string) bool {
	equal := true
	if len(tokens) != len(e.Tokens) {
		return false
	}
	for i, v := range e.Tokens {
		if v != tokens[i] {
			if contains(e.IgnoreDiffs, i) {
				log.WithFields(log.Fields{
					"index":    i,
					"expected": v,
					"actual":   tokens[i],
					"allowed":  true,
				}).Debug("deviate")
			} else {
				log.WithFields(log.Fields{
					"index":    i,
					"expected": v,
					"actual":   tokens[i],
					"allowed":  false,
				}).Debug("deviate")
				equal = false
			}
		}
	}
	return equal
}

// Diff builds the index set of differences between e.Tokens and tokens.
func (e Expectation) Diff(tokens []string) ([]int, error) {
	if len(tokens) != len(e.Tokens) {
		return []int{}, errors.New("number of tokes must be equal")
	}

	var diffs []int
	for i, v := range tokens {
		if v != e.Tokens[i] {
			diffs = append(diffs, i)
		}
	}
	return diffs, nil
}

func contains[T comparable](values []T, value T) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func (e Expectation) String() string {
	return strings.Join(e.Tokens, " ")
}

func (e Expectation) Shorten(i int) string {
	if i >= len(e.Tokens) {
		return e.String()
	}

	return fmt.Sprintf("%s...%s", strings.Join(e.Tokens[0:i/2], " "), strings.Join(e.Tokens[len(e.Tokens)-i/2:], " "))
}
