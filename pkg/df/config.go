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
	return config
}

func exists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
