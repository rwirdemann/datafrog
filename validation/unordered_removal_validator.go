package validation

import (
	"fmt"
	"github.com/rwirdemann/texttools/matcher"
	"log"
	"strings"
)

type UnorderedRemovalValidator struct {
	expectations []string
}

func NewUnorderedRemovalValidator(expectations []string) Validator {
	v := UnorderedRemovalValidator{}
	for _, e := range expectations {
		if len(strings.Trim(e, " ")) > 0 {
			v.expectations = append(v.expectations, e)
		}
	}
	return &v
}

func (u *UnorderedRemovalValidator) RemoveFirstMatchingExpectation(pattern string) {
	for i, expectation := range u.expectations {
		p := matcher.NewPattern(pattern)
		if p.MatchesInclude(expectation) && !p.MatchesExclude(expectation) {
			u.expectations = append(u.expectations[:i], u.expectations[i+1:]...)
			return
		}
	}
	log.Fatalf("Could not remove expectation. Pattern not found: %s", pattern)
}

func (u *UnorderedRemovalValidator) PrintResults() {
	if len(u.expectations) == 0 {
		fmt.Printf("All expectations met!")
	} else {
		fmt.Printf("Failed due to unmet expectations! Missing: %d", len(u.expectations))
	}
}
