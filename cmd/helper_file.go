// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/straubt1/tfx/output"
)

func isFile(filename string) bool {
	f, err := os.Stat(filename)
	if err != nil || f == nil {
		return false
	}

	if f.IsDir() {
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
// (optional) append a new folder structure depth
// Return absolute path
func getDirectory(directory string, additional ...string) (string, error) {
	if directory != "" {
		if !isDirectory(directory) {
			return "", errors.New("directory is not valid")
		}
	} else {
		if directory == "" {
			output.Get().Message("Directory not supplied, creating a temp directory")
			tempDir, err := os.MkdirTemp("", "tfx-")
			if err != nil {
				return "", err
			}
			directory = tempDir
		}
	}

	// add additional path, this may not exist yet but we have verified the top directory does
	for _, a := range additional {
		directory = filepath.Join(directory, a)
	}

	directory, err := filepath.Abs(directory)
	if err != nil {
		return "", errors.New("failed to get absolute directory")
	}

	return directory, nil
}
