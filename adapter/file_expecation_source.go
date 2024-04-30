package adapter

import (
	"bufio"
	"encoding/json"
	"github.com/rwirdemann/datafrog/app/domain"
	"log"
	"os"
)

// A FileExpectationSource reads expectations from a file in json format.
type FileExpectationSource struct {
	filename string
	testcase domain.Testcase
}

// NewFileExpectationSource creates a new NewFileExpectationSource that reads
// its expectations from filename.
func NewFileExpectationSource(filename string) (*FileExpectationSource, error) {
	expectations, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	fes := FileExpectationSource{filename: filename}
	if err := json.Unmarshal(expectations, &fes.testcase); err != nil {
		return nil, err
	}
	return &fes, nil
}

// Get returns the testcase.
func (s FileExpectationSource) Get() domain.Testcase {
	return s.testcase
}

func (s FileExpectationSource) Write(testcase domain.Testcase) error {
	f, err := os.Create(s.filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	b, err := json.Marshal(testcase)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := writer.WriteString(string(b)); err != nil {
		return err
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	return nil
}
