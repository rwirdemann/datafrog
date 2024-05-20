package ports

import "time"

type Timer interface {
	Start()
	GetStart() time.Time
	MatchesRecordingPeriod(ts time.Time) bool
}
