package df

import "time"

// Log defines methods to obtain the next line from the monitored log file and
// to extract its timestamp.
type Log interface {
	Timestamp(s string) (time.Time, error)
	NextLine() (string, error)
	Close()
	Tail() error
}
