package utils

import (
	"io/ioutil"
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

func GetDirFileNames(dirPath string) []string {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil || len(files) == 0 {
		return []string{}
	}
	result := make([]string, len(files))
	for i, file := range files {
		result[i] = file.Name()
	}
	return result
}
