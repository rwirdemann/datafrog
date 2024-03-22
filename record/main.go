package main

import (
	"bufio"
	"fmt"
	"github.com/rwirdemann/texttools/config"
	matcher2 "github.com/rwirdemann/texttools/matcher"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	c := config.NewConfig("config.json")
	fmt.Print("Press Enter to start recording...")
	_, _ = fmt.Scanln()

	recordingStartedAt := time.Now()
	fmt.Printf("Recording started at  %v. Press Enter to stop recording...\n", recordingStartedAt)

	readFile, _ := os.Open(c.Filename)
	defer func(readFile *os.File) {
		err := readFile.Close()
		if err != nil {
			log.Fatal(err)

		}
	}(readFile)

	go checkExit()

	out, err := os.Create(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {

		}
	}(out)
	outWriter := bufio.NewWriter(out)

	matcher := matcher2.NewPatternMatcher(c)
	reader := bufio.NewReader(readFile)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				time.Sleep(500 * time.Millisecond)
				continue
			}

			break
		}
		ts, validTimestamp := containsValidTimestamp(line)
		if validTimestamp && matchesRecordingPeriod(ts, recordingStartedAt) && matcher.MatchesAny(line) {
			fmt.Print(line)
			_, err := outWriter.WriteString(line)
			if err != nil {
				log.Fatal(err)
			}
			err = outWriter.Flush()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func checkExit() {
	var b = make([]byte, 1)
	l, _ := os.Stdin.Read(b)
	if l > 0 {
		os.Exit(0)
	}
}
func matchesRecordingPeriod(ts time.Time, startDate time.Time) bool {
	return ts.Equal(startDate) || ts.After(startDate)
}

func containsValidTimestamp(line string) (time.Time, bool) {
	split := strings.Split(line, "\t")
	if len(split) == 0 {
		return time.Time{}, false
	}

	d, err := time.Parse(time.RFC3339Nano, split[0])
	if err != nil {
		return time.Time{}, false
	}
	return d, true

}
