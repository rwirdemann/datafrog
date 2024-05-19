package cmd

import (
	"github.com/rwirdemann/datafrog/adapter"
	"github.com/rwirdemann/datafrog/internal/datafrog"
	"github.com/rwirdemann/datafrog/ports"
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
		c := datafrog.NewConfig("config.json")
		databaseLog := adapter.NewPostgresLog(c.Filename, c)
		defer databaseLog.Close()
		return Start(databaseLog)
	},
}
