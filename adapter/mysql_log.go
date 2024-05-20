package adapter

import (
	"bufio"
	"errors"
	"github.com/rwirdemann/datafrog/pkg/df"
	"io"
	"log"
	"os"
	"time"
)

type MySQLLog struct {
	logfile *os.File
	reader  *bufio.Reader
}

func NewMYSQLLog(logfileName string) df.Log {
	logfile, err := os.Open(logfileName)
	if err != nil {
		log.Fatal(err)
	}
	return MySQLLog{logfile: logfile, reader: bufio.NewReader(logfile)}
}

func (m MySQLLog) Close() {
	err := m.logfile.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (m MySQLLog) Timestamp(s string) (time.Time, error) {
	t, err := df.Timestamp(s, "[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\\.[0-9]{6}Z", time.RFC3339Nano)
	if err != nil {
		return time.Time{}, errors.New("string contains no valid Timestamp")
	}
	return t, nil
}

func (m MySQLLog) NextLine() (string, error) {
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
