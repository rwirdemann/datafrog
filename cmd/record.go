package cmd

import (
	"fmt"
	"github.com/rwirdemann/databasedragon/app"
	"log"
	"os"

	"github.com/rwirdemann/databasedragon/adapter"
	"github.com/rwirdemann/databasedragon/config"
	"github.com/rwirdemann/databasedragon/matcher"
	"github.com/spf13/cobra"
)

func init() {
	recordCmd.Flags().String("out", "", "Filename to save recording")
	recordCmd.Flags().Bool("prompt", false, "Wait for key stroke before recording starts")
	_ = recordCmd.MarkFlagRequired("out")
	rootCmd.AddCommand(recordCmd)
}

// close done channel to stop recording loop.
var recordingDone = make(chan struct{})

// read from stopped channel to wait for the recorder to finish
var recordingStopped = make(chan struct{})

var recorder *app.Recorder
var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Starts recording",
	Run: func(cmd *cobra.Command, args []string) {
		out, _ := cmd.Flags().GetString("out")
		c := config.NewConfig("config.json")
		prompt, _ := cmd.Flags().GetBool("prompt")
		if prompt {
			log.Printf("Recording goes to '%s'. Hit enter when you are ready!", out)
			_, _ = fmt.Scanln()
		} else {
			log.Printf("Recording goes to '%s'.", out)
		}

		recordingSink := adapter.NewFileRecordingSink(out)
		databaseLog := createLogAdapter(c)
		t := &adapter.UTCTimer{}
		recorder = app.NewRecorder(c, matcher.MySQLTokenizer{}, databaseLog, recordingSink, t)
		go checkExit()
		go recorder.Start(recordingDone, recordingStopped)
		<-recordingStopped
	},
}

// Checks if enter was hit to stop recording.
func checkExit() {
	var b = make([]byte, 1)
	l, _ := os.Stdin.Read(b)
	if l > 0 {
		close(recordingDone)
	}
}
