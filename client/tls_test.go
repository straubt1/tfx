// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package client

import (
	"net/http"
	"testing"
)

func TestHTTPClientForTFE_SkipVerify(t *testing.T) {
	t.Run("disabled by default", func(t *testing.T) {
		c := HTTPClientForTFE(false)
		if c.Transport != nil {
			t.Fatal("expected nil transport when ssl skip verify is disabled")
		}
	})

	t.Run("enabled sets InsecureSkipVerify", func(t *testing.T) {
		c := HTTPClientForTFE(true)
		tr, ok := c.Transport.(*http.Transport)
		if !ok {
			t.Fatalf("expected *http.Transport, got %T", c.Transport)
		}
		if tr.TLSClientConfig == nil || !tr.TLSClientConfig.InsecureSkipVerify {
			t.Fatal("expected InsecureSkipVerify to be true")
		}
	})
}

func TestNewHTTPClientWithLogging_SkipVerify(t *testing.T) {
	client, closer, err := NewHTTPClientWithLogging(nil, true)
	if err != nil {
		t.Fatalf("NewHTTPClientWithLogging() error = %v", err)
	}
	if closer != nil {
		defer func() { _ = closer.Close() }()
	}

	lt, ok := client.Transport.(*LoggingTransport)
	if !ok {
		t.Fatalf("expected *LoggingTransport, got %T", client.Transport)
	}
	tr, ok := lt.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("expected *http.Transport, got %T", lt.Transport)
	}
	if tr.TLSClientConfig == nil || !tr.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("expected InsecureSkipVerify to be true on logging transport")
	}
}
