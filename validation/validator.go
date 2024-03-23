package validation

type Validator interface {
	RemoveFirstMatchingExpectation(pattern string)
	PrintResults()
}
