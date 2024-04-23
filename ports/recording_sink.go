package ports

// RecordingSink defines methods to write recorded expectations.
type RecordingSink interface {
	WriteString(s string) (int, error)
	Flush() error
	Close()
}
