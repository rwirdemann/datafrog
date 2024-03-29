package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/rwirdemann/databasedragon/config"
	"github.com/spf13/cobra"
)

func init() {
	recordCmd.Flags().String("out", "", "Filename to save recording")
	recordCmd.MarkFlagRequired("out")
	rootCmd.AddCommand(recordCmd)
}

var recorder *Recorder
var recordCmd = &cobra.Command{
	Use:   "dbd",
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
