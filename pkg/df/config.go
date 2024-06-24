// Package df provides usecase overlapping datatypes and functions, like
// Expectation, Testcase or Tokenizer. Underlying design guideline: if a type
// belongs to an usecase it should live in the usecase package. If a type or
// function is used by multiple usecase it should live in df.
package df

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Config represents the settings used for record and verification runs.
// Patterns specify a list of patterns a log statement must match in order to be
// recorded or verified. Example: "select job!publish_trials<1" contains an in-
// and exclude rule thus only statements that contain `select job` but not
// `publish_trials<1` are considered.
type Config struct {
	SUT struct {
		BaseURL string `json:"base_url"` // base URL of SUT
	} `json:"sut"`
	Channels     []Channel `json:"channels"` // list of monitored channels
	Expectations struct {
		// report additional expectations that are not port of the initial
		// recording run
		ReportAdditional bool `json:"report_additional"`
	}
	// which ui driver: Playwright | none
	UIDriver   string `json:"ui_driver"`
	Playwright struct {
		BaseDir string `json:"base_dir"` // base directory of playwright project
		TestDir string `json:"test_dir"` // subdirectory in BaseDir where the tests are stored
	}
	Web struct {
		Port    int `json:"port"`    // web app http port
		Timeout int `json:"timeout"` // http timeout in seconds
	}

	Api struct {
		Port int `json:"port"` // api http port
	}
}

// NewDefaultConfig creates a new Config instance by trying to find a file
// config.json in the current or in the config subdirectory.
func NewDefaultConfig() (Config, error) {
	if exists("config.json") {
		return NewConfig("config.json"), nil
	}
	if exists("config/config.json") {
		return NewConfig("config/config.json"), nil
	}

	return Config{}, errors.New("config.json not found")
}

// NewConfig creates a new instance given its settings from filename in json
// format.
func NewConfig(filename string) Config {
	log.Printf("using config file '%s'", filename)
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
	if config.Channels[0].Format == "postgres" {
		var logFilePath = config.Channels[0].Log

		if strings.Contains(logFilePath, "YYYY-MM-DD-hhmmss.log") {

			// Logfile name contains a timestamp: Find the newest one in containing folder
			path := filepath.Dir(logFilePath)
			var expectedFileNameStart string = logFilePath[len(path) : len(logFilePath)-21]
			entries, err := os.ReadDir(path)
			if err == nil {
				var newestTime int64 = 0
				for _, file := range entries {
					// First check, if the name matches
					if strings.Contains(file.Name(), expectedFileNameStart) {
						// Now check, if it is the newest
						fi, err := os.Stat(path + file.Name())
						if err == nil {
							currTime := fi.ModTime().Unix()
							if currTime > newestTime {
								newestTime = currTime
								logFilePath = path + file.Name()
							}
						}
					}
				}
				// Set the resolved log file path top config
				config.Channels[0].Log = logFilePath
			}
		}
	}
	return config
}

func exists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
