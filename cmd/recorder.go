package cmd

import (
	"bufio"
	"log"
	"os"

	"github.com/rwirdemann/databasedragon/adapter"
	"github.com/rwirdemann/databasedragon/config"
	"github.com/rwirdemann/databasedragon/matcher"
	"github.com/rwirdemann/databasedragon/ticker"
)

type Recorder struct {
	config      config.Config
	outFilename string
	running     bool
}

func NewRecorder(c config.Config, outFilename string) *Recorder {
	return &Recorder{config: c, outFilename: outFilename, running: false}
}

// Start starts the recording process as endless loop. Every log entry that matches one of the
// patterns specified in config is written to the out file. Only log entries that fall in the actual
// recording period are considered. The caller should stop the recording by calling Recorder.Stop().
func (r *Recorder) Start() {
	r.running = true
	t := ticker.Ticker{}
	t.Start()
	log.Printf("Recording started at %v. Press Enter to stop recording...", t.GetStart())
	logPort := adapter.NewMYSQLLog(r.config.Filename)
	defer logPort.Close()

	out, err := os.Create(r.outFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	outWriter := bufio.NewWriter(out)

	m := matcher.NewDynamicDataMatcher(r.config)
	for {
		if !r.running {
			break
		}
		line, err := logPort.NextLine()
		if err != nil {
			log.Fatal(err)
		}

		ts, err := logPort.Timestamp(line)
		if err != nil {
			continue
		}
		if t.MatchesRecordingPeriod(ts) && m.MatchesAny(line) {
			log.Println(line)
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

// Stop stops the recording.
func (r *Recorder) Stop() {
	r.running = false
	log.Println("Recording stoped!")
}
