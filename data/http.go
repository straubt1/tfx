// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	"io"
	"net/http"
	"os"
)

// UploadBinary performs a PUT of the file at path to the given pre-signed URL
func UploadBinary(uploadURL string, path string) error {
	data, err := os.Open(path)
	if err != nil {
		return err
	}
	defer data.Close()

	req, err := http.NewRequest("PUT", uploadURL, data)
	if err != nil {
		return err
	}

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

// DownloadTextFile fetches the content at downloadURL and returns it as a string
func DownloadTextFile(downloadURL string) (string, error) {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	resp, err := client.Get(downloadURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
