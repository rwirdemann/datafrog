package postgres

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/rwirdemann/datafrog/pkg/df"
)

type Log struct {
	logfile *os.File
	reader  *bufio.Reader
}

func NewPostgresLog(logfileName string) Log {
	logfile, err := os.Open(logfileName)
	if err != nil {
		log.Fatal(err)
	}
	return Log{logfile: logfile, reader: bufio.NewReader(logfile)}
}

func (m Log) Close() {
	err := m.logfile.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (m Log) Timestamp(s string) (time.Time, error) {
	t, err := df.Timestamp(s, "[0-9]{4}-[0-9]{2}-[0-9]{2}\\s[0-9]{2}:[0-9]{2}:[0-9]{2}\\.[0-9]{3}", time.DateTime)
	if err != nil {
		return time.Time{}, errors.New("string contains no valid Timestamp")
	}
	return t, nil
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

func (m Log) NextLine(done chan struct{}) (string, error) {
	for {
		select {
		default:
			line, err := m.reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					time.Sleep(500 * time.Millisecond)
					continue
				}
				return "", err
			}

			matches := strings.Contains(line, "$1")
			if matches {
				line = m.mergeNext(line)
			}

			return line, nil
		case <-done:
			log.Printf("nextline: done channel closed")
			return "", nil
		}
	}
}

func (m Log) mergeNext(line string) string {
	next, _ := m.reader.ReadString('\n')
	if strings.Contains(next, "DETAIL") {
		values := make(map[int]string)
		r := regexp.MustCompile(`\$\d\s=\s'(?:[^']|'')*'|\$\d\s=\sNULL`)
		matches := r.FindAllStringSubmatch(next, -1)
		for i, v := range matches {
			split := strings.Split(v[0], "=")
			values[i] = strings.Trim(split[len(split)-1], " ")
		}
		i := 1
		for {
			placeholder := fmt.Sprintf("$%d", i)
			if strings.Contains(line, placeholder) {
				line = strings.Replace(line, placeholder, values[i-1], -1)
				i = i + 1
			} else {
				break
			}
		}
	}
	return line
}
