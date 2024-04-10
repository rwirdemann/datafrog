package adapter

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"

	"github.com/rwirdemann/databasedragon/ports"
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
	return Timestamp(s)
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
