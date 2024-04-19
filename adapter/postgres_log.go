package adapter

import (
	"bufio"
	"fmt"
	"github.com/rwirdemann/databasedragon/config"
	"github.com/rwirdemann/databasedragon/matcher"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

type PostgresLog struct {
	logfile *os.File
	reader  *bufio.Reader
	config  config.Config
}

func NewPostgresLog(logfileName string, config config.Config) PostgresLog {
	logfile, err := os.Open(logfileName)
	if err != nil {
		log.Fatal(err)
	}
	return PostgresLog{logfile: logfile, reader: bufio.NewReader(logfile), config: config}
}

func (m PostgresLog) Close() {
	err := m.logfile.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (m PostgresLog) Timestamp(s string) (time.Time, error) {
	return Timestamp(s)
}

func (m PostgresLog) NextLine() (string, error) {
	for {
		line, err := m.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			return "", err
		}

		matches, _ := matcher.MatchesPattern(m.config, line)
		if matches {
			next, _ := m.reader.ReadString('\n')
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

		return line, nil
	}
}
