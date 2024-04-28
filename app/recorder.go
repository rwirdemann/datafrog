package app

import (
	"encoding/json"
	"github.com/rwirdemann/databasedragon/app/domain"
	"github.com/rwirdemann/databasedragon/config"
	"github.com/rwirdemann/databasedragon/matcher"
	"github.com/rwirdemann/databasedragon/ports"
	"log"
)

// A Recorder monitors a database log and records all statements that match one
// of the patterns specified in config. The recorded output is written to
// recording sink.
type Recorder struct {
	config        config.Config
	tokenizer     matcher.Tokenizer
	databaseLog   ports.Log
	recordingSink ports.RecordingSink
	timer         ports.Timer
}

// NewRecorder creates a new Recorder.
func NewRecorder(c config.Config, tokenizer matcher.Tokenizer, log ports.Log, sink ports.RecordingSink, timer ports.Timer) *Recorder {
	return &Recorder{config: c, tokenizer: tokenizer, databaseLog: log, recordingSink: sink, timer: timer}
}

// Start starts the recording process as endless loop. Every log entry that
// matches one of the patterns specified in config is written to the recording
// sink. Only log entries that fall in the actual recording period are
// considered.
func (r *Recorder) Start(done chan struct{}, stopped chan struct{}) {
	r.timer.Start()
	log.Printf("Recording started at %v. Press Enter to stop and save recording...", r.timer.GetStart())
	var expectations []domain.Expectation

	// tell caller that verification has been finished
	defer close(stopped)

	// called when done channel is closed
	defer func() {
		r.writeExpectations(expectations)
		r.recordingSink.Close()
		r.databaseLog.Close()
	}()

	for {
		select {
		default:
			line, err := r.databaseLog.NextLine()
			if err != nil {
				log.Fatal(err)
			}

			ts, err := r.databaseLog.Timestamp(line)
			if err != nil {
				continue
			}
			if r.timer.MatchesRecordingPeriod(ts) {
				matches, pattern := matcher.MatchesPattern(r.config, line)
				if matches {
					tokens := r.tokenizer.Tokenize(line, r.config.Patterns)
					e := domain.Expectation{Tokens: tokens, IgnoreDiffs: []int{}, Pattern: pattern}
					expectations = append(expectations, e)
					log.Printf("new expectation: %s\n", e.Shorten(8))
				}
			}
		case <-done:
			log.Println("Recording finished. Run verification now!")
			return
		}
	}
}

// writeExpectations writes initialExpectations as json to the recordingSink.
// Existing exceptions are overridden.
func (r *Recorder) writeExpectations(expectations []domain.Expectation) {
	b, err := json.Marshal(expectations)
	if err != nil {
		log.Fatal(err)
	}
	_, err = r.recordingSink.WriteString(string(b))
	if err != nil {
		log.Fatal(err)
	}
	err = r.recordingSink.Flush()
	if err != nil {
		log.Fatal(err)
	}
}
