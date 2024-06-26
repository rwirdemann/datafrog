package mysql

import "github.com/rwirdemann/datafrog/pkg/df"

type LogFactory struct {
}

func (f LogFactory) Create(filename string) (df.Log, error) {
	log, err := NewMYSQLLog(filename)
	return log, err
}
