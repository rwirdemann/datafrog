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

// The Verifier verifies the expectations in expectationSource. It monitors the
// databsaseLog for these expectations and requires them to be in same order as
// given in expectationSource. Verified expectations are written to
// verificationSink.
type Verifier struct {
	config            config.Config
	databsaseLog      ports.Log
	expectationSource ports.ExpectationSource
	verificationSink  ports.RecordingSink
	timer             ports.Timer
	running           bool
}

// NewVerifier creates a new Verifier.
func NewVerifier(c config.Config, log ports.Log, source ports.ExpectationSource, sink ports.RecordingSink, t ports.Timer) *Verifier {
	return &Verifier{
		config:            c,
		databsaseLog:      log,
		expectationSource: source,
		verificationSink:  sink,
		timer:             t,
		running:           false}
}

// Start runs the verification loop. It stops, when the expectations got out of
// order or when all expectations where met.
func (v *Verifier) Start() error {
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

				// Since we expect the expectation in order we always remove the
				// first from the list. But only if it matches the current
				// matching pattern. If it doesn't match we return an error
				// because the verify run didn't receive the expectations in the
				// required order.
				if err := v.expectationSource.RemoveFirst(pattern); err != nil {
					return err
				}

				log.Printf("Verfication met: '%s'", line)
				_, err := v.verificationSink.WriteString(line)
				if err != nil {
					return err
				}

				err = v.verificationSink.Flush()
				if err != nil {
					return err
				}

				if len(v.expectationSource.GetAll()) == 0 {
					v.Stop()
					break
				}
			}
		}
	}
	return nil
}

// Stop stops the verifcation.
func (v *Verifier) Stop() {
	v.running = false
	if len(v.expectationSource.GetAll()) == 0 {
		log.Println("Verfication was successful!")
	} else {
		log.Printf("Verfication failed. %d unmatched verfications\n", len(v.expectationSource.GetAll()))
	}
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
		return verifier.Start()
	},
}

// Checks if enter was hit to stop verfication.
func checkVerifyExit() {
	var b = make([]byte, 1)
	l, _ := os.Stdin.Read(b)
	if l > 0 {
		verifier.Stop()
		os.Exit(0)
	}
}
