package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Config struct {
	Filename string `json:"filename"`
}

func NewConfig() Config {
	configfile, err := os.Open("config.json")
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
