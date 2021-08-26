package utils

import (
	"os"
)

func IsPathExist(path string) bool {
	// create video dictionary if needed
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true // may be other reason to cause error
}
