package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rwirdemann/datafrog/pkg/df"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

type TestRepository struct{}

func (r TestRepository) All() ([]df.Testcase, error) {
	var all []df.Testcase
	dir, err := os.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("TestRepository.All failed: %w", err)
	}
	for _, f := range dir {
		if strings.HasSuffix(f.Name(), ".json") && !strings.HasPrefix(f.Name(), "config") {
			tc, err := r.Get(f.Name())
			if err != nil {
				if errors.Is(err, InvalidJsonError{}) {
					log.Errorf("testfile '%s' contains invalid json. deleting file.", f.Name())
					_ = os.Remove(f.Name())
				} else {
					log.Errorf("TestRepository.Get failed: %v", err)
				}
			} else {
				all = append(all, tc)
			}
		}
	}
	return all, nil
}

type InvalidJsonError struct{}

func (e InvalidJsonError) Error() string {
	return "json: invalid data"
}

func (r TestRepository) Get(filename string) (df.Testcase, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return df.Testcase{}, err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.Debugf("error closing file %s: %v", filename, err)
		}
	}(jsonFile)
	b, _ := io.ReadAll(jsonFile)
	if len(b) == 0 {
		return df.Testcase{}, InvalidJsonError{}
	}
	var tc df.Testcase
	if err := json.Unmarshal(b, &tc); err != nil {
		return df.Testcase{}, InvalidJsonError{}
	}
	return tc, nil
}
