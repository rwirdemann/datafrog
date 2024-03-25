package ticker

import "time"

type Ticker struct {
	start time.Time
}

func (t *Ticker) Start() {
	t.start = time.Now()
}

func (t *Ticker) GetStart() time.Time {
	return t.start
}

func (t *Ticker) MatchesRecordingPeriod(ts time.Time) bool {
	return ts.Equal(t.start) || ts.After(t.start)
}
