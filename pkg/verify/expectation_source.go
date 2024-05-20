package verify

import (
	"github.com/rwirdemann/datafrog/pkg/df"
)

// ExpectationSource defines methods to read and write expectations from and to
// an underlying source.
type ExpectationSource interface {
	Get() df.Testcase
	Write(testcase df.Testcase) error
}
