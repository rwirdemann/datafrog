package cmd

import (
	"github.com/rwirdemann/databasedragon/adapter"
	"github.com/rwirdemann/databasedragon/config"
	"github.com/rwirdemann/databasedragon/ports"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	rootCmd.AddCommand(tailCmd)
}

func Start(l ports.Log) error {
	for {
		line, err := l.NextLine()
		if err != nil {
			log.Fatal(err)
		}
		println(line)
	}
	return nil
}

var tailCmd = &cobra.Command{
	Use:   "tail",
	Short: "Tails log file",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := config.NewConfig("config.json")
		databaseLog := adapter.NewPostgresLog(c.Filename, c)
		defer databaseLog.Close()
		return Start(databaseLog)
	},
}
