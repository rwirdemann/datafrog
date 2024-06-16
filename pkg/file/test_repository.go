package file

import (
	"encoding/json"
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
				log.Errorf("TestRepository.Get failed: %v", err)
			} else {
				all = append(all, tc)
			}
		}
	}
	return all, nil
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
		return df.Testcase{}, fmt.Errorf("testfile '%s' contains no data", filename)
	}
	var tc df.Testcase
	if err := json.Unmarshal(b, &tc); err != nil {
		return df.Testcase{}, err
	}
	return tc, nil
}
