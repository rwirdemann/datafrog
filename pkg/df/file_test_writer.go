package df

import (
	log "github.com/sirupsen/logrus"
	"os"
)

type FileTestWriter struct {
	file     *os.File
	filename string
}

func NewFileTestWriter(filename string) (*FileTestWriter, error) {
	ftw := &FileTestWriter{filename: filename}
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
	if err := f.file.Close(); err != nil {
		return err
	}
	log.Printf("%s closed", f.filename)
	return nil
}
