package record

import (
	"encoding/json"
	"github.com/rwirdemann/datafrog/internal/datafrog"
	"github.com/rwirdemann/datafrog/ports"
	"log"
)

// A Recorder monitors a database log and records all statements that match one
// of the patterns specified in config. The recorded output is written to
// recording sink.
type Recorder struct {
	config        datafrog.Config
	tokenizer     datafrog.Tokenizer
	log           ports.Log
	recordingSink ports.RecordingSink
	timer         ports.Timer
	name          string
	uuidProvider  ports.UUIDProvider
	testcase      datafrog.Testcase
}

// NewRecorder creates a new Recorder.
func NewRecorder(c datafrog.Config, tokenizer datafrog.Tokenizer,
	log ports.Log, sink ports.RecordingSink, timer ports.Timer, name string,
	uuidProvider ports.UUIDProvider) *Recorder {

	return &Recorder{
		config:        c,
		tokenizer:     tokenizer,
		log:           log,
		recordingSink: sink,
		timer:         timer,
		name:          name,
		uuidProvider:  uuidProvider,
		testcase:      datafrog.Testcase{Name: name}}
}

// Start starts the recording process as endless loop. Every log entry that
// matches one of the patterns specified in config is written to the recording
// sink. Only log entries that fall in the actual recording period are
// considered.
func (r *Recorder) Start(done chan struct{}, stopped chan struct{}) {
	r.timer.Start()
	log.Printf("Recording started at %v. Press Enter to stop and save recording...", r.timer.GetStart())

	// tell caller that verification has been finished
	defer close(stopped)

	// called when done channel is closed
	defer func() {
		r.write(r.testcase)
		r.recordingSink.Close()
		r.log.Close()
	}()

	for {
		select {
		default:
			line, err := r.log.NextLine()
			if err != nil {
				log.Fatal(err)
			}

			ts, err := r.log.Timestamp(line)
			if err != nil {
				continue
			}
			if r.timer.MatchesRecordingPeriod(ts) {
				matches, pattern := datafrog.MatchesPattern(r.config, line)
				if matches {
					tokens := r.tokenizer.Tokenize(line, r.config.Patterns)
					e := datafrog.Expectation{Uuid: r.uuidProvider.NewString(), Tokens: tokens, IgnoreDiffs: []int{}, Pattern: pattern}
					r.testcase.Expectations = append(r.testcase.Expectations, e)
					log.Printf("new expectation: %s\n", e.Shorten(8))
				}
			}
		// check if the caller (web, cli, ...) has closed the done channel to
		// tell me that recoding has been finished
		case <-done:
			log.Println("Recording finished. Run verification now!")
			return
		}
	}
}

// write writes testcase as json to the recordingSink. Existing data is
// overridden.
func (r *Recorder) write(testcase datafrog.Testcase) {
	b, err := json.Marshal(testcase)
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

func (r *Recorder) Testcase() datafrog.Testcase {
	return r.testcase
}
