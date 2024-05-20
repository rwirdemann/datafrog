package df

import (
	"errors"
	"regexp"
	"time"
)

// Timestamp finds the first Timestamp in s that matches the pattern and returns
// a time.Time created by using the given layout.
func Timestamp(s, pattern, layout string) (time.Time, error) {
	ts := regexp.MustCompile(pattern).FindString(s)
	if d, err := time.Parse(layout, ts); err != nil {
		return time.Time{}, errors.New("string contains no valid Timestamp")
	} else {
		return d, nil
	}
}
