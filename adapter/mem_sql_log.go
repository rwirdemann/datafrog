package adapter

import (
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
	return Timestamp(s)
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
