package matcher

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Expectation struct {
	tokens      []string
	ignoreDiffs []int
}

func NewExpectation(expectation string, verification string) Expectation {
	e := Expectation{tokens: tokenize(expectation), ignoreDiffs: buildDiff(expectation, verification)}
	return e
}

func (e Expectation) Equal(actual string) bool {
	equal := true
	actualTokens := tokenize(actual)
	if len(actualTokens) != len(e.tokens) {
		return false
	}
	for i, v := range e.tokens {
		if v != actualTokens[i] {
			if contains(e.ignoreDiffs, i) {
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

func tokenize(s string) []string {
	tokens := []string{}
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
	return tokens
}

func buildDiff(expectation, verification string) []int {
	t1 := tokenize(expectation)
	t2 := tokenize(verification)
	diffs := []int{}
	for i, v := range t1 {
		if v != t2[i] {
			diffs = append(diffs, i)
		}
	}
	return diffs
}
