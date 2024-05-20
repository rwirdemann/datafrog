package file

import (
	"bufio"
	"log"
	"os"
)

type RecordingSink struct {
	outFile   *os.File
	outWriter *bufio.Writer
}

func NewFileRecordingSink(filename string) RecordingSink {
	sink := RecordingSink{}
	var err error
	sink.outFile, err = os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	sink.outWriter = bufio.NewWriter(sink.outFile)
	return sink
}

func (fs RecordingSink) WriteString(s string) (int, error) {
	return fs.outWriter.WriteString(s)
}

func (fs RecordingSink) Flush() error {
	return fs.outWriter.Flush()
}

func (fs RecordingSink) Close() {
	fs.outFile.Close()
}
