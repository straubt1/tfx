// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package file

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/straubt1/tfx/output"
)

// IsFile returns true if the given path is a file
func IsFile(filename string) bool {
	f, err := os.Stat(filename)
	if err != nil || f == nil {
		return false
	}

	if f.IsDir() {
		return false
	}

	return true
}

// IsDirectory returns true if the given path is a directory
func IsDirectory(filename string) bool {
	f, err := os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	if !f.IsDir() {
		return false
	}
	return true
}

// ReadFile reads the contents of a file and returns it as a string
func ReadFile(filename string) (string, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return "", nil
	}

	return string(b), nil
}

// GetDirectory validates or creates a directory, optionally appending additional path segments.
// Given a directory, validate it is a real directory.
// If no directory, create a temp directory.
// (optional) append a new folder structure depth.
// Return absolute path.
func GetDirectory(directory string, additional ...string) (string, error) {
	if directory != "" {
		if !IsDirectory(directory) {
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

	if err := os.MkdirAll(directory, 0755); err != nil {
		return "", errors.New("failed to create directory")
	}

	return directory, nil
}
