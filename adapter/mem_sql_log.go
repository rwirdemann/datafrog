package adapter

import (
	"errors"
	"strings"
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
	split := strings.Split(s, "\t")
	if len(split) == 0 {
		return time.Time{}, errors.New("string contains no valid timestamp")
	}

	d, err := time.Parse(time.RFC3339Nano, split[0])
	if err != nil {
		return time.Time{}, errors.New("string contains no valid timestamp")
	}
	return d, nil
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
