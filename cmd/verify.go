package cmd

import (
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
	verifyCmd.Flags().String("expectations", "", "Filename to save verify")
	verifyCmd.MarkFlagRequired("expectations")
	rootCmd.AddCommand(verifyCmd)
}

type Verifier struct {
	config            config.Config
	databsaseLog      ports.Log
	expectationSource ports.ExpectationSource
	verificationSink  ports.RecordingSink
	timer             ports.Timer
	running           bool
}

func NewVerifier(c config.Config, log ports.Log, source ports.ExpectationSource, sink ports.RecordingSink, t ports.Timer) *Verifier {
	return &Verifier{
		config:            c,
		databsaseLog:      log,
		expectationSource: source,
		verificationSink:  sink,
		timer:             t,
		running:           false}
}

func (v *Verifier) Start() {
	v.running = true
	v.timer.Start()
	log.Printf("Verification started at %v. Press Enter to stop verification...", v.timer.GetStart())

	for {
		if !v.running {
			break
		}
		line, err := v.databsaseLog.NextLine()
		if err != nil {
			log.Fatal(err)
		}

		// Hack to enable test adapter to stop the recording
		if line == "STOP" {
			break
		}

		ts, err := v.databsaseLog.Timestamp(line)
		if err != nil {
			continue
		}
		if v.timer.MatchesRecordingPeriod(ts) {
			matches, pattern := matcher.MatchesPattern(v.config, line)
			if matches {
				log.Println(line)
				_, err := v.verificationSink.WriteString(line)
				if err != nil {
					log.Fatal(err)
				}

				err = v.verificationSink.Flush()
				if err != nil {
					log.Fatal(err)
				}

				if err := v.expectationSource.RemoveFirst(pattern); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

// Stop stops the verifcation.
func (v *Verifier) Stop() {
	v.running = false
	log.Println("Recording stoped!")
}

var verifier *Verifier
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Starts verifcation",
	RunE: func(cmd *cobra.Command, args []string) error {
		expectationsFilename, _ := cmd.Flags().GetString("expectations")
		verficationFilename := fmt.Sprintf("%s.verify", expectationsFilename)
		c := config.NewConfig("config.json")
		log.Printf("Verifying '%s'. Verification goes to '%s'. Hit enter when you are ready!", expectationsFilename, verficationFilename)
		_, _ = fmt.Scanln()
		go checkVerifyExit()

		expectationSource := adapter.NewFileExpectationSource(expectationsFilename)
		verificationSink := adapter.NewFileRecordingSink(verficationFilename)
		defer verificationSink.Close()

		databaseLog := adapter.NewMYSQLLog(c.Filename)
		defer databaseLog.Close()

		t := &adapter.UTCTimer{}

		verifier = NewVerifier(c, databaseLog, expectationSource, verificationSink, t)
		recorder.Start()

		return nil
	},
}

// Checks if enter was hit to stop recording.
func checkVerifyExit() {
	var b = make([]byte, 1)
	l, _ := os.Stdin.Read(b)
	if l > 0 {
		verifier.Stop()
		os.Exit(0)
	}
}
