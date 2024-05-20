package mocks

import "time"

type Timer struct {
	start time.Time
}

func (t Timer) Start() {
}

func (t Timer) GetStart() time.Time {
	return t.start
}

func (t Timer) MatchesRecordingPeriod(ts time.Time) bool {
	return true
}
