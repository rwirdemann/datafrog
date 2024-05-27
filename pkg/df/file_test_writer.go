package df

import (
	"os"
)

type FileTestWriter struct {
	file *os.File
}

func NewFileTestWriter(filename string) (*FileTestWriter, error) {
	ftw := &FileTestWriter{}
	var err error
	ftw.file, err = os.Create(filename)
	if err != nil {
		return nil, err
	}
	return ftw, nil
}

func (f FileTestWriter) Write(p []byte) (int, error) {
	return f.file.Write(p)
}

func (f FileTestWriter) Close() error {
	return f.file.Close()
}
