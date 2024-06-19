package df

import "time"

type UTCTimer struct {
	start time.Time
}

func (t *UTCTimer) Start() {
	t.start = time.Now().UTC()
}

func (t *UTCTimer) GetStart() time.Time {
	return t.start
}

func (t *UTCTimer) MatchesRecordingPeriod(ts time.Time) bool {
	return ts.Equal(t.start) || ts.After(t.start)
}
