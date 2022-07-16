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

func readFile(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", nil
	}

	return string(b), nil
}
