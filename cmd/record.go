package cmd

import (
	"bufio"
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
	recordCmd.MarkFlagRequired("out")
	rootCmd.AddCommand(recordCmd)
}

type Recorder struct {
	config       config.Config
	databsaseLog ports.Log
	timer        ports.Timer
	outFilename  string
	running      bool
}

func NewRecorder(c config.Config, databsaseLog ports.Log, timer ports.Timer, outFilename string) *Recorder {
	return &Recorder{config: c, databsaseLog: databsaseLog, timer: timer, outFilename: outFilename, running: false}
}

// Start starts the recording process as endless loop. Every log entry that matches one of the
// patterns specified in config is written to t he out file. Only log entries that fall in the
// actual recording period are considered. The caller should stop the recording by calling
// Recorder.Stop().
func (r *Recorder) Start() {
	r.running = true
	r.timer.Start()
	log.Printf("Recording started at %v. Press Enter to stop recording...", t.GetStart())

	out, err := os.Create(r.outFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	outWriter := bufio.NewWriter(out)

	for {
		if !r.running {
			break
		}
		line, err := r.databsaseLog.NextLine()
		if err != nil {
			log.Fatal(err)
		}

		ts, err := r.databsaseLog.Timestamp(line)
		if err != nil {
			continue
		}
		if r.timer.MatchesRecordingPeriod(ts) && matcher.MatchesPattern(r.config, line) {
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

		databaseLog := adapter.NewMYSQLLog(c.Filename)
		defer databaseLog.Close()

		t := &adapter.UTCTimer{}

		recorder = NewRecorder(c, databaseLog, t, out)
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
