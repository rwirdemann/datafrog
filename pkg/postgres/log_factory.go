package postgres

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwirdemann/datafrog/pkg/df"
)

type LogFactory struct {
}

func (f LogFactory) Create(filename string) df.Log {
	logFilePath, err := resolveDate(filename)
	if err != nil {
		log.Fatalf("LogFactory: Could not resolve Log-File %s: %s", filename, err)
	}
	return NewPostgresLog(logFilePath)
}

func resolveDate(filename string) (string, error) {

	var dateString = "YYYY-MM-DD_hhmmss.log"
	if !strings.Contains(filename, dateString) {
		return filename, nil
	}

	// Logfile name contains a timestamp: Find the newest one in containing folder

	path := filepath.Dir(filename)

	entries, err := os.ReadDir(path)
	if err != nil {
		return filename, err
	}

	var expectedFileNameStart string = filename[len(path)+1 : len(filename)-len(dateString)]
	var logFilePath = filename
	var newestTime int64 = 0
	var found = false
	for _, file := range entries {
		// First check, if the name matches
		if strings.Contains(file.Name(), expectedFileNameStart) {
			// Now check, if it is the newest
			fi, err := os.Stat(path + "/" + file.Name())
			if err == nil {
				currTime := fi.ModTime().Unix()
				if currTime > newestTime {
					found = true
					newestTime = currTime
					logFilePath = path + "/" + file.Name()
				}
			}
		}
	}

	err = nil
	if !found {
		err = fs.ErrNotExist
	}
	return logFilePath, err
}
