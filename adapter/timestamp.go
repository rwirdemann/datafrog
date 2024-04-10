package adapter

import (
	"errors"
	"strings"
	"time"
)

func Timestamp(s string) (time.Time, error) {
	s = strings.ReplaceAll(s, "\t", " ")
	split := strings.Split(s, " ")
	if len(split) == 0 {
		return time.Time{}, errors.New("string contains no valid timestamp")
	}

	d, err := time.Parse(time.RFC3339Nano, split[0])
	if err != nil {
		return time.Time{}, errors.New("string contains no valid timestamp")
	}
	return d, nil
}
