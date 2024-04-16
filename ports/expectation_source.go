package ports

import "github.com/rwirdemann/databasedragon/matcher"

// ExpectationSource defines methods to read and write expectations from and to
// an underlying source.
type ExpectationSource interface {
	GetAll() []matcher.Expectation
	WriteAll()
}
