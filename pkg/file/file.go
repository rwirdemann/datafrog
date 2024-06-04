package file

import "os"

// Exists returns true if filename exists and false otherwise.
func Exists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
