package main

import (
	"fmt"
	"strings"
)

func main() {
	expectation := "insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World', '2024-04-02 08:37:37', 0, null, '', 'Hello', 39)"
	verification := "insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World', '2024-04-02 08:37:38', 0, null, '', 'Hello', 40)"
	actual := "insert into job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World', '2024-04-02 08:37:39', 0, null, '', 'Hello', 41)"
	actual2 := "update job (description, publish_at, publish_trials, published_timestamp, tags, title, id) values ('World', '2024-04-02 08:37:39', 0, null, '', 'Hello', 41)"

	diffs := buildDiff(expectation, verification)
	compare(expectation, actual, diffs)
	compare(expectation, actual2, diffs)
}

func compare(expectation, actual string, diffs []int) {
	t1 := tokenize(expectation)
	t2 := tokenize(actual)
	for i, v := range t1 {
		if v != t2[i] {
			fmt.Printf("Token %d varies: %s / %s...", i, v, t2[i])
			if contains(diffs, i) {
				fmt.Printf("thats OK\n")
			} else {
				fmt.Printf("thats NOT OK\n")
			}
		}
	}
}

func contains(values []int, value int) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func buildDiff(expectation, verification string) []int {
	t1 := tokenize(expectation)
	t2 := tokenize(verification)
	diffs := []int{}
	for i, v := range t1 {
		if v != t2[i] {
			fmt.Printf("Token %d varies: %s / %s\n", i, v, t2[i])
			diffs = append(diffs, i)
		}
	}
	return diffs
}

func tokenize(s string) []string {
	split := strings.Split(s, ",")
	var tokens = []string{}
	for _, v := range split {
		tokens = append(tokens, strings.Trim(v, " )"))
	}

	return tokens
}
