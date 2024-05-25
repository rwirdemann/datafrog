package mysql

import (
	"bufio"
	"errors"
	"github.com/rwirdemann/datafrog/pkg/df"
	"io"
	"log"
	"os"
	"time"
)

type Log struct {
	logfile *os.File
	reader  *bufio.Reader
}

func NewMYSQLLog(logfileName string) Log {
	logfile, err := os.Open(logfileName)
	if err != nil {
		log.Fatal(err)
	}
	return Log{logfile: logfile, reader: bufio.NewReader(logfile)}
}

// Tail sets the read cursor of the log file to its end.
func (m Log) Tail() error {
	log.Printf("tailing %s...", m.logfile.Name())
	defer log.Printf("tailing successful!")
	for {
		_, err := m.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		}
	}
}

func (m Log) Close() {
	err := m.logfile.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (m Log) Timestamp(s string) (time.Time, error) {
	t, err := df.Timestamp(s, "[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\\.[0-9]{6}Z", time.RFC3339Nano)
	if err != nil {
		return time.Time{}, errors.New("string contains no valid Timestamp")
	}
	return t, nil
}

func (m Log) NextLine() (string, error) {
	for {
		line, err := m.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			return "", err
		}
		return line, nil
	}
}
