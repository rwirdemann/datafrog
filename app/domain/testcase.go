package domain

type Testcase struct {
	Name          string        `json:"name"`
	Running       bool          `json:"running"`
	Verifications int           `json:"verifications"`
	Expectations  []Expectation `json:"expectation"`

	// Expectations, that match one of the patterns but didn't match one of the
	// expected expectations
	AdditionalExpectations []Expectation `json:"additional_expectations"`
}
