package postgres

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rwirdemann/datafrog/pkg/df"
)

type LogFactory struct {
}

func (f LogFactory) Create(logFilePath string) df.Log {
	var dateString = "YYYY-MM-DD_hhmmss.log"
	if strings.Contains(logFilePath, dateString) {

		// Logfile name contains a timestamp: Find the newest one in containing folder
		path := filepath.Dir(logFilePath)
		var expectedFileNameStart string = logFilePath[len(path)+1 : len(logFilePath)-len(dateString)]
		entries, err := os.ReadDir(path)
		if err == nil {
			var newestTime int64 = 0
			for _, file := range entries {
				// First check, if the name matches
				if strings.Contains(file.Name(), expectedFileNameStart) {
					// Now check, if it is the newest
					fi, err := os.Stat(path + "/" + file.Name())
					if err == nil {
						currTime := fi.ModTime().Unix()
						if currTime > newestTime {
							newestTime = currTime
							logFilePath = path + "/" + file.Name()
						}
					}
				}
			}
		}
	}
	return NewPostgresLog(logFilePath)
}
