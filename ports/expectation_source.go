package ports

import (
	"github.com/rwirdemann/databasedragon/app/domain"
)

// ExpectationSource defines methods to read and write expectations from and to
// an underlying source.
type ExpectationSource interface {
	Get() domain.Testcase
	Write(testcase domain.Testcase) error
}
