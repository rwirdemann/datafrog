package validation

type Validator interface {
	RemoveFirstMatchingExpectation(pattern string)
	Remove(expectation string)
	GetExpectations() []string
	PrintResults()
}
