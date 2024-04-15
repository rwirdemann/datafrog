package adapter

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/rwirdemann/databasedragon/matcher"
)

// A FileExpectationSource reads expectations from a file in json format.
type FileExpectationSource struct {
	expectations []matcher.Expectation
	filename     string
}

// NewFileExpectationSource creates a new NewFileExpectationSource that reads
// its expectations from filename.
func NewFileExpectationSource(filename string) *FileExpectationSource {
	expectations, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	fes := FileExpectationSource{filename: filename}
	if err := json.Unmarshal(expectations, &fes.expectations); err != nil {
		log.Fatal(err)
	}

	return &fes
}

// GetAll returns all expectations.
func (s *FileExpectationSource) GetAll() []matcher.Expectation {
	return s.expectations
}

func (s *FileExpectationSource) WriteAll([]matcher.Expectation) {
	f, err := os.Create(s.filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	b, err := json.Marshal(s.expectations)
	if err != nil {
		log.Fatal(err)
	}
	writer.WriteString(string(b))
	writer.Flush()
}
