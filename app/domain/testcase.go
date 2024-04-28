package domain

type Testcase struct {
	Name          string        `json:"name"`
	Running       bool          `json:"running"`
	Verifications int           `json:"verifications"`
	Expectations  []Expectation `json:"expectation"`
}
