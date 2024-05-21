package mocks

type MemWriter struct {
	Recorded []string
}

func (ms *MemWriter) Write(p []byte) (n int, err error) {
	ms.Recorded = append(ms.Recorded, string(p))
	return len(p), nil
}
