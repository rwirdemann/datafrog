package ports

type RecordingSink interface {
	WriteString(s string) (int, error)
	Flush() error
	Close()
}
