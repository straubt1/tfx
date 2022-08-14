package cmd

import (
	"errors"
	"io/ioutil"
	"os"
)

func isFile(filename string) bool {
	if f, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) || f.IsDir() {
		return false
	}
	return true
}

func isDirectory(filename string) bool {
	f, err := os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	if !f.IsDir() {
		return false
	}
	return true
}

func readFile(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", nil
	}

	return string(b), nil
}

// Given a directory, validate it is a real directory
// If no directory, create a temp directory
// Return path
func getDirectory(directory string) (string, error) {
	if directory != "" {
		if !isDirectory(directory) {
			return "", errors.New("directory is not valid")
		}
	} else {
		o.AddMessageUserProvided("Directory not supplied, creating a temp directory", "")
		dst, err := ioutil.TempDir("", "slug")
		if err != nil {
			return "", errors.New("failed to create temp directory")
		}
		directory = dst
	}

	return directory, nil
}
