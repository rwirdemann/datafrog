package main

import (
	"bufio"
	"fmt"
	"github.com/rwirdemann/texttools/config"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var (
	patterns = []string{"insert into job", "update job"}
)

var expectations []string

func main() {
	config := config.NewConfig("config.json")
	println(config.Filename)
	fmt.Print("Press Enter to start listening...")
	_, _ = fmt.Scanln()

	listeningStartedAt := time.Now()
	fmt.Printf("Listening started at  %v. Press Enter to stop recording...\n", listeningStartedAt)

	validationFile, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	split := strings.Split(string(validationFile), "\n")
	for _, s := range split {
		if len(strings.Trim(s, " ")) > 0 {
			expectations = append(expectations, s)
		}
	}

	readFile, _ := os.Open(config.Filename)
	defer readFile.Close()

	go checkExit()

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
		if validTimestamp && matchesRecordingPeriod(ts, listeningStartedAt) {
			pattern, matches := matchesPattern(line)
			if matches {
				fmt.Printf("Expectation met: %s", line)
				expectations = remove(expectations, pattern)
			}
		}
	}
}

func remove(expectations []string, pattern string) []string {
	var result []string
	for i, expectation := range expectations {
		if strings.Contains(expectation, pattern) {
			result = append(expectations[:i], expectations[i+1:]...)
			return result
		}
	}
	return expectations
}

func checkExit() {
	var b = make([]byte, 1)
	l, _ := os.Stdin.Read(b)
	if l > 0 {
		if len(expectations) == 0 {
			fmt.Printf("All expectations met!")
		} else {
			fmt.Printf("Failed due to unmet expectations! Missing: %d", len(expectations))
		}
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

func matchesPattern(line string) (string, bool) {
	for _, pattern := range patterns {
		if strings.Contains(line, pattern) {
			return pattern, true
		}
	}
	return "", false
}
