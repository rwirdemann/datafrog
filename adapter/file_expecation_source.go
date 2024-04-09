package adapter

import (
	"log"
	"os"
	"strings"
)

type FileExpectationSource struct {
	filename string
}

func NewFileExpectationSource(filename string) FileExpectationSource {
	return FileExpectationSource{filename: filename}
}

func (s FileExpectationSource) GetAll() []string {
	expectations, err := os.ReadFile(s.filename)
	if err != nil {
		log.Fatal(err)
	}
	split := strings.Split(string(expectations), "\n")
	a := []string{}
	for _, v := range split {
		if strings.Trim(v, " ") != "" {
			a = append(a, v)
		}
	}
	return a
}

// RemoveFirst is not implemented for file based source.
func (s FileExpectationSource) RemoveFirst(pattern string) error {
	return nil
}
