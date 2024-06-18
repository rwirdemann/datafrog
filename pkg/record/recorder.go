package record

import (
	"context"
	"encoding/json"
	"github.com/rwirdemann/datafrog/pkg/df"
	"log"
)

// A Recorder monitors a database log and records all statements that match one
// of the patterns specified in config. The recorded output is written to
// recording sink.
type Recorder struct {
	config       df.Config
	tokenizer    df.Tokenizer
	log          df.Log
	writer       df.TestWriter // destination of recorded testcase
	timer        df.Timer
	name         string
	uuidProvider UUIDProvider
	testcase     df.Testcase
}

// NewRecorder creates a new Recorder.
func NewRecorder(c df.Config, tokenizer df.Tokenizer,
	log df.Log, w df.TestWriter, timer df.Timer, name string,
	uuidProvider UUIDProvider) *Recorder {

	return &Recorder{
		config:       c,
		tokenizer:    tokenizer,
		log:          log,
		writer:       w,
		timer:        timer,
		name:         name,
		uuidProvider: uuidProvider,
		testcase:     df.Testcase{Name: name}}
}

// Start starts the recording process as endless loop. Every log entry that
// matches one of the patterns specified in config is written to the recording
// sink. Only log entries that fall in the actual recording period are
// considered.
func (r *Recorder) Start(done chan struct{}, stopped chan struct{}) {
	r.timer.Start()
	log.Printf("Recording started at %v...", r.timer.GetStart())

	// tell caller that verification has been finished
	defer close(stopped)

	// called when done channel is closed
	defer func() {
		b, err := json.Marshal(r.testcase)
		if err != nil {
			log.Fatal(err)
		}
		_, err = r.writer.Write(b)
		if err != nil {
			log.Fatal(err)
		}
		if err := r.writer.Close(); err != nil {
			log.Fatal(err)
		}
		r.log.Close()
	}()

	// jump to log file end
	if err := r.log.Tail(); err != nil {
		log.Fatal(err)
	}

	// create a cancel context to interrupt blocking NextLine call when this function
	// terminated
	cancelContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		select {
		default:
			line, err := r.log.NextLine(cancelContext)
			if err != nil {
				log.Fatal(err)
			}

			ts, err := r.log.Timestamp(line)
			if err != nil {
				continue
			}
			if r.timer.MatchesRecordingPeriod(ts) {
				matches, pattern := df.MatchesPattern(r.config.Patterns, line)
				if matches {
					tokens := r.tokenizer.Tokenize(line, r.config.Patterns)
					e := df.Expectation{Uuid: r.uuidProvider.NewString(), Tokens: tokens, IgnoreDiffs: []int{}, Pattern: pattern}
					r.testcase.Expectations = append(r.testcase.Expectations, e)
					log.Printf("new expectation: %s\n", e.Shorten(8))
				}
			}
		// check if the caller (web, cli, ...) has closed the done channel to
		// tell me that recoding has been finished
		case <-done:
			log.Println("recorder: done channel closed")
			return
		}
	}
}

func (r *Recorder) Testcase() df.Testcase {
	return r.testcase
}
