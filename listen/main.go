package main

import (
	"bufio"
	"fmt"
	"github.com/rwirdemann/texttools/config"
	matcher2 "github.com/rwirdemann/texttools/matcher"
	"github.com/rwirdemann/texttools/validation"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var validator validation.Validator

func main() {
	c := config.NewConfig("config.json")
	println(c.Filename)
	fmt.Print("Press Enter to start listening...")
	_, _ = fmt.Scanln()

	listeningStartedAt := time.Now()
	fmt.Printf("Listening started at  %v. Press Enter to stop recording...\n", listeningStartedAt)

	validationFile, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	validator = validation.NewUnorderedRemovalValidator(strings.Split(string(validationFile), "\n"))

	logFile, _ := os.Open(c.Filename)
	defer func(readFile *os.File) {
		err := readFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(logFile)

	go checkExit()

	matcher := matcher2.NewPatternMatcher(c)
	reader := bufio.NewReader(logFile)
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
		if validTimestamp && matchesRecordingPeriod(ts, listeningStartedAt) {
			if matcher.MatchesAny(line) {
				fmt.Printf("Expectation met: %s", line)
				pattern := matcher.MatchingPattern(line)
				validator.RemoveFirstMatchingExpectation(pattern)
			}
		}
	}
}

func checkExit() {
	var b = make([]byte, 1)
	l, _ := os.Stdin.Read(b)
	if l > 0 {
		validator.PrintResults()
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
