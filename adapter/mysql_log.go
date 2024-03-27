package adapter

import (
	"bufio"
	"errors"
	"github.com/rwirdemann/databasedragon/ports"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type MySQLLog struct {
	logfile *os.File
	reader  *bufio.Reader
}

func NewMYSQLLog(logfileName string) ports.Log {
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
	split := strings.Split(s, "\t")
	if len(split) == 0 {
		return time.Time{}, errors.New("string contains no valid timestamp")
	}

	d, err := time.Parse(time.RFC3339Nano, split[0])
	if err != nil {
		return time.Time{}, errors.New("string contains no valid timestamp")
	}
	return d, nil
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
