package adapter

import (
	"bufio"
	"encoding/json"
	"github.com/rwirdemann/databasedragon/app/domain"
	"log"
	"os"
)

// A FileExpectationSource reads expectations from a file in json format.
type FileExpectationSource struct {
	expectations []domain.Expectation
	filename     string
}

// NewFileExpectationSource creates a new NewFileExpectationSource that reads
// its expectations from filename.
func NewFileExpectationSource(filename string) (*FileExpectationSource, error) {
	expectations, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	fes := FileExpectationSource{filename: filename}
	if err := json.Unmarshal(expectations, &fes.expectations); err != nil {
		return nil, err
	}

	return &fes, nil
}

// GetAll returns all expectations.
func (s *FileExpectationSource) GetAll() []domain.Expectation {
	return s.expectations
}

func (s *FileExpectationSource) WriteAll() {
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
