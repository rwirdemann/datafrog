package df

import (
	"context"
	"time"
)

// Log defines methods to obtain the next line from the monitored log file and
// to extract its timestamp.
type Log interface {
	Timestamp(s string) (time.Time, error)
	NextLine(ctx context.Context) (string, error)
	Close()
	Tail() error
}
