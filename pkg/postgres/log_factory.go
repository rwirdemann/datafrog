package postgres

import (
	"github.com/rwirdemann/datafrog/pkg/df"
)

type LogFactory struct {
}

func (f LogFactory) Create(filename string) df.Log {
	return NewPostgresLog(filename)
}
