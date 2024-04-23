package cmd

import (
	"github.com/rwirdemann/databasedragon/adapter"
	"github.com/rwirdemann/databasedragon/config"
	"github.com/rwirdemann/databasedragon/ports"
	"log"
)

func createLogAdapter(c config.Config) ports.Log {
	switch c.Logformat {
	case "mysql":
		return adapter.NewMYSQLLog(c.Filename)
	case "postgres":
		return adapter.NewPostgresLog(c.Filename, c)
	default:
		log.Fatalf("Unknown log format: %s", c.Logformat)
	}
	return nil
}
