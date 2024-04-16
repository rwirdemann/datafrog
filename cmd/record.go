package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/rwirdemann/databasedragon/adapter"
	"github.com/rwirdemann/databasedragon/config"
	"github.com/rwirdemann/databasedragon/matcher"
	"github.com/rwirdemann/databasedragon/ports"
	"github.com/spf13/cobra"
)

func init() {
	recordCmd.Flags().String("out", "", "Filename to save recording")
	recordCmd.Flags().Bool("prompt", false, "Wait for key stroke before recording starts")
	recordCmd.MarkFlagRequired("out")
	rootCmd.AddCommand(recordCmd)
}

// A Recorder monitors a database log and records all statements that match one
// of the patterns specified in config. The recorded output is written to
// recording sink.
type Recorder struct {
	config        config.Config
	databsaseLog  ports.Log
	recordingSink ports.RecordingSink
	timer         ports.Timer
	running       bool
}

// NewRecorder creates a new Recorder.
func NewRecorder(c config.Config, log ports.Log, sink ports.RecordingSink, timer ports.Timer) *Recorder {
	return &Recorder{config: c, databsaseLog: log, recordingSink: sink, timer: timer, running: false}
}

// Start starts the recording process as endless loop. Every log entry that
// matches one of the patterns specified in config is written to the recording
// sink. Only log entries that fall in the actual recording period are
// considered. The caller should stop the recording by calling Recorder.Stop().
func (r *Recorder) Start() {
	r.running = true
	r.timer.Start()
	log.Printf("Recording started at %v. Press Enter to stop recording...", r.timer.GetStart())
	var expectations []matcher.Expectation
	for {
		if !r.running {
			r.writeExpectation(expectations)
			break
		}
		line, err := r.databsaseLog.NextLine()
		if err != nil {
			log.Fatal(err)
		}

		// Hack to enable test adapter to stop the recording
		if line == "STOP" {
			r.writeExpectation(expectations)
			break
		}

		ts, err := r.databsaseLog.Timestamp(line)
		if err != nil {
			continue
		}
		if r.timer.MatchesRecordingPeriod(ts) {
			matches, pattern := matcher.MatchesPattern(r.config, line)
			if matches {
				log.Println(line)
				tokens := matcher.Tokenize(matcher.Normalize(line, r.config.Patterns))
				e := matcher.Expectation{Tokens: tokens, IgnoreDiffs: []int{}, Pattern: pattern}
				expectations = append(expectations, e)
			}
		}
	}
}

func (r *Recorder) writeExpectation(expectations []matcher.Expectation) {
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

// Stop stops the recording.
func (r *Recorder) Stop() {
	r.running = false
	log.Println("Recording finshed. Run verfication now!")
}

var recorder *Recorder
var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Starts recording",
	RunE: func(cmd *cobra.Command, args []string) error {
		out, _ := cmd.Flags().GetString("out")
		c := config.NewConfig("config.json")
		prompt, _ := cmd.Flags().GetBool("prompt")
		if prompt {
			log.Printf("Recording goes to '%s'. Hit enter when you are ready!", out)
			_, _ = fmt.Scanln()
		} else {
			log.Printf("Recording goes to '%s'.", out)
		}
		go checkExit()

		recordingSink := adapter.NewFileRecordingSink(out)
		defer recordingSink.Close()

		databaseLog := adapter.NewMYSQLLog(c.Filename)
		defer databaseLog.Close()

		t := &adapter.UTCTimer{}

		recorder = NewRecorder(c, databaseLog, recordingSink, t)
		recorder.Start()

		return nil
	},
}

// Checks if enter was hit to stop recording.
func checkExit() {
	var b = make([]byte, 1)
	l, _ := os.Stdin.Read(b)
	if l > 0 {
		recorder.Stop()
	}
}
