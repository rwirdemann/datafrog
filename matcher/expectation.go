package matcher

import (
	"fmt"
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
	for i, v := range e.tokens {
		if v != actualTokens[i] {
			fmt.Printf("Token %d differs: %s / %s...", i, v, actualTokens[i])
			if contains(e.ignoreDiffs, i) {
				fmt.Printf("thats OK\n")
			} else {
				fmt.Printf("thats NOT OK\n")
				equal = false
			}
		}
	}
	return equal
}

func contains(values []int, value int) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func tokenize(s string) []string {
	tokens := []string{}
	t := ""
	quoted := false
	for i := 0; i < len(s); i++ {
		if string(s[i]) == "'" {
			quoted = !quoted
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
