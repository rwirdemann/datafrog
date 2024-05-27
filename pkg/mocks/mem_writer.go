package mocks

import (
	"encoding/json"
	"github.com/rwirdemann/datafrog/pkg/df"
)

type MemWriter struct {
	Testcase df.Testcase
}

func (ms *MemWriter) Close() error {
	return nil
}

func (ms *MemWriter) Write(p []byte) (n int, err error) {
	if err := json.Unmarshal(p, &ms.Testcase); err != nil {
		return 0, err
	}
	return len(p), nil
}
