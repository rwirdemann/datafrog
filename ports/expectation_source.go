package ports

import (
	"github.com/rwirdemann/datafrog/internal/datafrog"
)

// ExpectationSource defines methods to read and write expectations from and to
// an underlying source.
type ExpectationSource interface {
	Get() datafrog.Testcase
	Write(testcase datafrog.Testcase) error
}
