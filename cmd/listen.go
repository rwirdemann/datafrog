package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/rwirdemann/databasedragon/config"
	"github.com/spf13/cobra"
)

func init() {
	listenCmd.Flags().String("expectations", "", "Filename with expectations")
	listenCmd.MarkFlagRequired("expectations")
	rootCmd.AddCommand(listenCmd)
}

var listener *Listener
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Starts listening and validation",
	RunE: func(cmd *cobra.Command, args []string) error {
		expectations, _ := cmd.Flags().GetString("expectations")
		c := config.NewConfig("config.json")
		log.Printf("Listening to '%s'. Hit enter when you are ready!", expectations)
		_, _ = fmt.Scanln()
		go checkStopListening()

		listener = NewListener(c, expectations)
		listener.Start()

		return nil
	},
}

// Checks if enter was hit to stop listening.
func checkStopListening() {
	var b = make([]byte, 1)
	l, _ := os.Stdin.Read(b)
	if l > 0 {
		listener.Stop()
		os.Exit(0)
	}
}
