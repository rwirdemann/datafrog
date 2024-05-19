package datafrog

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

// Config represents the settings used for record and verification runs.
// Patterns specify a list of patterns a log statement must match in order to be
// recorded or verified. Example: "select job!publish_trials<1" contains an in-
// and exclude rule thus only statements that contain `select job` but not
// `publish_trials<1` are considered.
type Config struct {
	Filename     string   `json:"filename"`  // full path of logfile to be used
	Logformat    string   `json:"logformat"` // format of log file, actual mysql | postgresql
	Patterns     []string `json:"patterns"`  // list of patterns to consider
	Expectations struct {
		// report additional expectations that are not port of the initial
		// recording run
		ReportAdditional bool `json:"report_additional"`
	}
	Web struct {
		Port    int `json:"port"`    // web app http port
		Timeout int `json:"timeout"` // http timeout in seconds
	}

	Api struct {
		Port int `json:"port"` // api http port
	}
}

// NewConfig creates a new instance given its settings from filename in json
// format.
func NewConfig(filename string) Config {
	configfile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func(configfile *os.File) {
		err := configfile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(configfile)
	byteValue, _ := io.ReadAll(configfile)
	var config Config
	if err := json.Unmarshal(byteValue, &config); err != nil {
		log.Fatal(err)
	}
	return config
}
