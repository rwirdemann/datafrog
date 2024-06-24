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

type JSONTestRepository struct{}

func (r JSONTestRepository) Delete(testname string) error {
	if err := os.Remove(fmt.Sprintf("%s.json", testname)); err != nil {
		return err
	}
	return nil
}

func (r JSONTestRepository) Write(testname string, testcase df.Testcase) error {
	f, err := os.Create(fmt.Sprintf("%s.json", testname))
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			log.Errorf("unable to close file: %v", err)
		}
	}(f)
	if err != nil {
		return err
	}
	b, err := json.Marshal(testcase)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write(b)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("successfully wrote %s\n", f.Name())
	return nil
}

func (r JSONTestRepository) Exists(testname string) bool {
	fn := testname
	if !strings.HasSuffix(testname, ".json") {
		fn = fmt.Sprintf("%s.json", testname)
	}

	if _, err := os.Stat(fn); os.IsNotExist(err) {
		return false
	}
	return true
}

func (r JSONTestRepository) All() ([]df.Testcase, error) {
	var all []df.Testcase
	dir, err := os.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("JSONTestRepository.All failed: %w", err)
	}
	for _, f := range dir {
		if strings.HasSuffix(f.Name(), ".json") && !strings.HasPrefix(f.Name(), "config") {
			tc, err := r.Get(f.Name())
			if err != nil {
				if errors.Is(err, InvalidJsonError{}) {
					log.Errorf("testfile '%s' contains invalid json. deleting file.", f.Name())
					_ = os.Remove(f.Name())
				} else {
					log.Errorf("JSONTestRepository.Get failed: %v", err)
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

func (r JSONTestRepository) Get(testname string) (df.Testcase, error) {
	fn := testname
	if !strings.HasSuffix(testname, ".json") {
		fn = fmt.Sprintf("%s.json", testname)
	}
	f, err := os.Open(fn)
	if err != nil {
		return df.Testcase{}, err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.Debugf("JSONTestRepository: error closing file %s: %v", fn, err)
		}
	}(f)
	b, _ := io.ReadAll(f)
	if len(b) == 0 {
		return df.Testcase{}, InvalidJsonError{}
	}
	var tc df.Testcase
	if err := json.Unmarshal(b, &tc); err != nil {
		return df.Testcase{}, InvalidJsonError{}
	}
	return tc, nil
}
