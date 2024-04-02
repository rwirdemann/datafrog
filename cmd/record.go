package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/rwirdemann/databasedragon/adapter"
	"github.com/rwirdemann/databasedragon/config"
	"github.com/rwirdemann/databasedragon/matcher"
	"github.com/rwirdemann/databasedragon/ticker"
	"github.com/spf13/cobra"
)

func init() {
	recordCmd.Flags().String("out", "", "Filename to save recording")
	recordCmd.MarkFlagRequired("out")
	rootCmd.AddCommand(recordCmd)
}

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

	m := matcher.NewLevenshteinMatcher(r.config)
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
		matches, _ := m.MatchesPattern(line)
		if t.MatchesRecordingPeriod(ts) && matches {
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

var recorder *Recorder
var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Starts recording",
	RunE: func(cmd *cobra.Command, args []string) error {
		out, _ := cmd.Flags().GetString("out")
		c := config.NewConfig("config.json")
		log.Printf("Recording goes to '%s'. Hit enter when you are ready!", out)
		_, _ = fmt.Scanln()
		go checkExit()

		recorder = NewRecorder(c, out)
		recorder.Start()

		return nil
	},
}

// Checks if enter was hit to stop recording
func checkExit() {
	var b = make([]byte, 1)
	l, _ := os.Stdin.Read(b)
	if l > 0 {
		recorder.Stop()
		os.Exit(0)
	}
}
