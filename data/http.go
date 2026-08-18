// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	"io"
	"net/http"
	"os"

	"github.com/pkg/errors"
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

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("download failed: %s returned %d", downloadURL, resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// DownloadFile streams the content at downloadURL to destPath. It fails on
// non-200 responses so a failed CDN fetch cannot be mistaken for a valid file.
func DownloadFile(downloadURL, destPath string) error {
	resp, err := http.Get(downloadURL)
	if err != nil {
		return errors.Wrap(err, "download request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("download failed: %s returned %d", downloadURL, resp.StatusCode)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return errors.Wrapf(err, "failed to create %s", destPath)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return errors.Wrapf(err, "failed to write %s", destPath)
	}
	return nil
}
