package cmd

import (
	"github.com/rwirdemann/datafrog/adapter"
	"github.com/rwirdemann/datafrog/internal/datafrog"
	"github.com/rwirdemann/datafrog/ports"
	"log"
)

func createLogAdapter(c datafrog.Config) ports.Log {
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
