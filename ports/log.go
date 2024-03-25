package ports

import "time"

type Log interface {
	Timestamp(s string) (time.Time, error)
	NextLine() (string, error)
	Close()
}
