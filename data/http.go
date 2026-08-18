// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/straubt1/tfx/version"
)

var publicRegistryHTTPClient = newPublicRegistryHTTPClient()

func newPublicRegistryHTTPClient() *http.Client {
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSHandshakeTimeout = 10 * time.Second
	tr.ResponseHeaderTimeout = 30 * time.Second
	return &http.Client{Transport: tr}
}

func publicRegistryUserAgent() string {
	return "tfx/" + version.Version
}

func publicRegistryGet(rawURL string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", publicRegistryUserAgent())
	return publicRegistryHTTPClient.Do(req)
}

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
	_, _ = io.Copy(io.Discard, res.Body)

	switch res.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		return nil
	default:
		return errors.Errorf("upload failed: %s returned %d", uploadURL, res.StatusCode)
	}
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
	resp, err := publicRegistryGet(downloadURL)
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
