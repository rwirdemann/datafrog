package adapter

import (
	"errors"
	"regexp"
	"time"
)

// timestamp finds the first timestamp in s that matches the pattern and returns
// a time.Time created by using the given layout.
func timestamp(s, pattern, layout string) (time.Time, error) {
	ts := regexp.MustCompile(pattern).FindString(s)
	if d, err := time.Parse(layout, ts); err != nil {
		return time.Time{}, errors.New("string contains no valid timestamp")
	} else {
		return d, nil
	}
}
