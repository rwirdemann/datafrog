package adapter

import (
	"bufio"
	"log"
	"os"
)

type FileRecordingSink struct {
	outFile   *os.File
	outWriter *bufio.Writer
}

func NewFileRecordingSink(filename string) FileRecordingSink {
	sink := FileRecordingSink{}
	var err error
	sink.outFile, err = os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	sink.outWriter = bufio.NewWriter(sink.outFile)
	return sink
}

func (fs FileRecordingSink) WriteString(s string) (int, error) {
	return fs.outWriter.WriteString(s)
}

func (fs FileRecordingSink) Flush() error {
	return fs.outWriter.Flush()
}

func (fs FileRecordingSink) Close() {
	fs.outFile.Close()
}
