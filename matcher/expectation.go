package matcher

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Expectation struct {
	Tokens      []string `json:"tokens"`
	Pattern     string
	Fulfilled   bool
	Verified    int
	IgnoreDiffs []int `json:"ignoreDiffs"`
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

func contains[T comparable](values []T, value T) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func (e Expectation) Diff(sql string) ([]int, error) {
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

func (e Expectation) String() string {
	return strings.Join(e.Tokens, " ")
}

func (e Expectation) Shorten(i int) string {
	if i >= len(e.Tokens) {
		return e.String()
	}

	return fmt.Sprintf("%s...%s", strings.Join(e.Tokens[0:i/2], " "), strings.Join(e.Tokens[len(e.Tokens)-i/2:], " "))
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
