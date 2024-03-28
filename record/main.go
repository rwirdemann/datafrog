package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rwirdemann/databasedragon/adapter"
	"github.com/rwirdemann/databasedragon/config"
	"github.com/rwirdemann/databasedragon/matcher"
	"github.com/rwirdemann/databasedragon/ticker"
)

func main() {
	c := config.NewConfig("config.json")
	fmt.Print("Press Enter to start recording...")
	_, _ = fmt.Scanln()

	t := ticker.Ticker{}
	t.Start()
	fmt.Printf("Recording started at  %v. Press Enter to stop recording...\n", t.GetStart())

	logPort := adapter.NewMYSQLLog(c.Filename)
	defer logPort.Close()

	go checkExit()

	out, err := os.Create(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	outWriter := bufio.NewWriter(out)

	m := matcher.NewDynamicDataMatcher(c)
	for {
		line, err := logPort.NextLine()
		if err != nil {
			log.Fatal(err)
		}

		ts, err := logPort.Timestamp(line)
		if err != nil {
			continue
		}
		if t.MatchesRecordingPeriod(ts) && m.MatchesAny(line) {
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
