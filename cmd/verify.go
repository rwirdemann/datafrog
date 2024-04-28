package cmd

import (
	"fmt"
	"github.com/rwirdemann/databasedragon/adapter"
	"github.com/rwirdemann/databasedragon/app"
	"github.com/rwirdemann/databasedragon/config"
	"github.com/rwirdemann/databasedragon/matcher"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func init() {
	verifyCmd.Flags().String("expectations", "", "Filename to save verify")
	verifyCmd.Flags().Bool("prompt", false, "Wait for key stroke before verification starts")
	verifyCmd.MarkFlagRequired("expectations")
	rootCmd.AddCommand(verifyCmd)
}

// close done channel to stop the verify loop.
var done = make(chan struct{})

// read from stopped channel to wait for the verifier to finish
var stopped = make(chan struct{})

var verifier *app.Verifier
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Starts verification",
	Run: func(cmd *cobra.Command, args []string) {
		expectationsFilename, _ := cmd.Flags().GetString("expectations")
		c := config.NewConfig("config.json")
		prompt, _ := cmd.Flags().GetBool("prompt")
		if prompt {
			log.Printf("Verifying '%s'. Hit enter when you are ready!", expectationsFilename)
			_, _ = fmt.Scanln()
		} else {
			log.Printf("Verifying '%s'.", expectationsFilename)
		}

		expectationSource, err := adapter.NewFileExpectationSource(expectationsFilename)
		if err != nil {
			log.Fatal(err)
		}
		databaseLog := createLogAdapter(c)
		defer databaseLog.Close()
		t := &adapter.UTCTimer{}
		verifier = app.NewVerifier(c, matcher.MySQLTokenizer{}, databaseLog, expectationSource, t, expectationsFilename)
		go checkVerifyExit()
		go verifier.Start(done, stopped)
		<-stopped // wait until verifier signals its finish
		fmt.Println(verifier.ReportResults())
	},
}

// checkVerifyExit if key  was hit to stop verification.
func checkVerifyExit() {
	var b = make([]byte, 1)
	l, _ := os.Stdin.Read(b)
	if l > 0 {
		close(done)
	}
}
