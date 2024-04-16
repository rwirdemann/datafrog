package matcher

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Expectation struct {
	Tokens      []string `json:"tokens"`
	Pattern     string
	Fulfilled   bool
	Verified    int
	IgnoreDiffs []int `json:"ignoreDiffs"`
}

func NewExpectation(expectation string, verification string) Expectation {
	diffs := buildDiff(expectation, verification)
	e := Expectation{Tokens: Tokenize(expectation), IgnoreDiffs: diffs}
	return e
}

func (e Expectation) Equal(actual string) bool {
	equal := true
	actualTokens := Tokenize(actual)
	if len(actualTokens) != len(e.Tokens) {
		return false
	}
	for i, v := range e.Tokens {
		if v != actualTokens[i] {
			if contains(e.IgnoreDiffs, i) {
				log.WithFields(log.Fields{
					"index":    i,
					"expected": v,
					"actual":   actualTokens[i],
					"allowed":  true,
				}).Debug("deviate")
			} else {
				log.WithFields(log.Fields{
					"index":    i,
					"expected": v,
					"actual":   actualTokens[i],
					"allowed":  false,
				}).Debug("deviate")
				equal = false
			}
		}
	}
	return equal
}

func (e Expectation) BuildDiff(sql string) ([]int, error) {
	tokens := Tokenize(sql)
	if len(tokens) != len(e.Tokens) {
		return []int{}, errors.New("number of tokes must be equals")
	}

	var diffs []int
	for i, v := range tokens {
		if v != e.Tokens[i] {
			diffs = append(diffs, i)
		}
	}
	return diffs, nil

}

func Tokenize(s string) []string {
	var tokens []string
	t := ""
	quoted := false
	for i := 0; i < len(s); i++ {
		if string(s[i]) == "'" {
			quoted = !quoted
			continue
		}

		if string(s[i]) != " " {
			t = fmt.Sprintf("%s%s", t, string(s[i]))
			continue
		}

		if string(s[i]) == " " && quoted {
			t = fmt.Sprintf("%s%s", t, string(s[i]))
			continue
		}

		if string(s[i]) == " " {
			tokens = append(tokens, t)
			t = ""
		}
	}
	tokens = append(tokens, t)
	log.Debug(tokens)
	return tokens
}

func buildDiff(expectation, verification string) []int {
	t1 := Tokenize(expectation)
	t2 := Tokenize(verification)
	diffs := []int{}
	for i, v := range t1 {
		if v != t2[i] {
			diffs = append(diffs, i)
		}
	}
	return diffs
}
