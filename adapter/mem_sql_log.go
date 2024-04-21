package adapter

import (
	"errors"
	"time"
)

type MemSQLLog struct {
	logs  []string
	index int
}

func NewMemSQLLog(logs []string) *MemSQLLog {
	return &MemSQLLog{logs: logs, index: 0}
}

func (l *MemSQLLog) Timestamp(s string) (time.Time, error) {
	t, err := timestamp(s, "[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\\.[0-9]{6}Z", time.RFC3339Nano)
	if err != nil {
		return time.Time{}, errors.New("string contains no valid timestamp")
	}
	return t, nil
}

func (l *MemSQLLog) NextLine() (string, error) {
	if l.index >= len(l.logs) {
		return "", nil
	}
	line := l.logs[l.index]
	l.index = l.index + 1
	return line, nil
}

func (l *MemSQLLog) Close() {

}
