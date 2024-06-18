package mocks

import (
	"errors"
	"github.com/rwirdemann/datafrog/pkg/df"
	"time"
)

type SQLLog struct {
	logs        []string
	index       int
	doneChannel chan struct{} // close this channel to notify verification loop to stop
}

func (l *SQLLog) Tail() error {
	return nil
}

func NewMemSQLLog(logs []string, doneChannel chan struct{}) *SQLLog {
	return &SQLLog{logs: logs, index: 0, doneChannel: doneChannel}
}

func (l *SQLLog) Timestamp(s string) (time.Time, error) {
	t, err := df.Timestamp(s, "[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\\.[0-9]{6}Z", time.RFC3339Nano)
	if err != nil {
		return time.Time{}, errors.New("string contains no valid timestamp")
	}
	return t, nil
}

func (l *SQLLog) NextLine(done chan struct{}) (string, error) {
	if l.index >= len(l.logs) {
		return "", nil
	}
	line := l.logs[l.index]
	l.index = l.index + 1
	if line == "STOP" {
		close(l.doneChannel)
	}
	return line, nil
}

func (l *SQLLog) Close() {

}
