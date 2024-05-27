package df

type TestWriter interface {
	Write(p []byte) (n int, err error)
	Close() error
}
