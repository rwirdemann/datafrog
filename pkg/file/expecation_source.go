package file

import (
	"bufio"
	"encoding/json"
	"github.com/rwirdemann/datafrog/pkg/df"
	"log"
	"os"
)

// A ExpectationSource reads expectations from a file in json format.
type ExpectationSource struct {
	filename string
	testcase df.Testcase
}

// NewFileExpectationSource creates a new NewFileExpectationSource that reads
// its expectations from filename.
func NewFileExpectationSource(filename string) (*ExpectationSource, error) {
	expectations, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	fes := ExpectationSource{filename: filename}
	if err := json.Unmarshal(expectations, &fes.testcase); err != nil {
		return nil, err
	}
	return &fes, nil
}

// Get returns the testcase.
func (s ExpectationSource) Get() df.Testcase {
	return s.testcase
}

func (s ExpectationSource) Write(testcase df.Testcase) error {
	f, err := os.Create(s.filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	// don't write additional expectations
	testcase.AdditionalExpectations = nil
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
