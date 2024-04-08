package adapter

import "time"

type MockTimer struct {
	start time.Time
}

func (t MockTimer) Start() {
}

func (t MockTimer) GetStart() time.Time {
	return t.start
}

func (t MockTimer) MatchesRecordingPeriod(ts time.Time) bool {
	return true
}
