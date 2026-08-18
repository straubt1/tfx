// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUploadBinaryAcceptsSuccessStatuses(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusCreated, http.StatusNoContent} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPut {
					t.Errorf("method = %s", r.Method)
				}
				_, _ = io.Copy(io.Discard, r.Body)
				w.WriteHeader(status)
			}))
			defer server.Close()

			path := filepath.Join(t.TempDir(), "file.bin")
			if err := os.WriteFile(path, []byte("payload"), 0644); err != nil {
				t.Fatal(err)
			}
			if err := UploadBinary(server.URL, path); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUploadBinaryRejectsErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		http.Error(w, "nope", http.StatusForbidden)
	}))
	defer server.Close()

	path := filepath.Join(t.TempDir(), "file.bin")
	if err := os.WriteFile(path, []byte("payload"), 0644); err != nil {
		t.Fatal(err)
	}
	err := UploadBinary(server.URL, path)
	if err == nil || !strings.Contains(err.Error(), "403") {
		t.Fatalf("expected 403 error, got %v", err)
	}
}

func TestPublicRegistryGetSetsUserAgent(t *testing.T) {
	var gotUA string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	resp, err := publicRegistryGet(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if !strings.HasPrefix(gotUA, "tfx/") {
		t.Fatalf("User-Agent = %q, want tfx/...", gotUA)
	}
}

func TestDownloadFileSetsUserAgent(t *testing.T) {
	var gotUA string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	dest := filepath.Join(t.TempDir(), "out.bin")
	if err := DownloadFile(server.URL, dest); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(gotUA, "tfx/") {
		t.Fatalf("User-Agent = %q, want tfx/...", gotUA)
	}
}
