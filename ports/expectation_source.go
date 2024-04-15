package ports

import "github.com/rwirdemann/databasedragon/matcher"

// ExpectationSource defines methods to read and write expectations from an
// underlying source.
type ExpectationSource interface {
	GetAll() []matcher.Expectation
	WriteAll(expectations []matcher.Expectation)
}
