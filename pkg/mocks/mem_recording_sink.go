package mocks

type RecordingSink struct {
	Recorded []string
}

func NewRecordingSink() *RecordingSink {
	return &RecordingSink{}
}

func (ms *RecordingSink) WriteString(s string) (int, error) {
	ms.Recorded = append(ms.Recorded, s)
	return 0, nil
}

func (ms *RecordingSink) Flush() error {
	return nil
}

func (ms *RecordingSink) Close() {
}
