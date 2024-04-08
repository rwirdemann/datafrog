package adapter

type MemRecordingSink struct {
	Recorded []string
}

func NewMemRecordingSink() *MemRecordingSink {
	return &MemRecordingSink{}
}

func (ms *MemRecordingSink) WriteString(s string) (int, error) {
	ms.Recorded = append(ms.Recorded, s)
	return 0, nil
}

func (ms *MemRecordingSink) Flush() error {
	return nil
}

func (ms *MemRecordingSink) Close() {
}
