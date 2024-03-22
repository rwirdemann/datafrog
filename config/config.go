package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Config struct {
	Filename string   `json:"filename"`
	Patterns []string `json:"patterns"`
}

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
