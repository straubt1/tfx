package client

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestRedactSensitiveData(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "redact authorization header",
			input:    "Authorization: Bearer token123\r\n",
			expected: "Authorization: [REDACTED]\r\n",
		},
		{
			name:     "redact cookie header",
			input:    "Cookie: session=abc123\r\n",
			expected: "Cookie: [REDACTED]\r\n",
		},
		{
			name:     "redact api key header",
			input:    "X-Api-Key: secretkey456\r\n",
			expected: "X-Api-Key: [REDACTED]\r\n",
		},
		{
			name:     "preserve non-sensitive headers",
			input:    "Content-Type: application/json\r\n",
			expected: "Content-Type: application/json\r\n",
		},
		{
			name:     "case insensitive redaction",
			input:    "authorization: Bearer token123\r\n",
			expected: "authorization: [REDACTED]\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := redactSensitiveData([]byte(tt.input))
			if string(result) != tt.expected {
				t.Errorf("redactSensitiveData() = %q, want %q", string(result), tt.expected)
			}
		})
	}
}

func TestLoggingTransportRoundTrip(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	}))
	defer server.Close()

	// Create a temporary log file
	tmpDir := t.TempDir()
	logFile, err := os.Create(filepath.Join(tmpDir, "test.log"))
	if err != nil {
		t.Fatal(err)
	}
	defer logFile.Close()

	// Create logging transport
	transport := &LoggingTransport{
		Transport: http.DefaultTransport,
		LogFile:   logFile,
	}

	// Create request
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Execute request
	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip() error = %v", err)
	}
	defer resp.Body.Close()

	// Verify response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Verify log file was written to
	logFile.Sync()
	stat, err := logFile.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if stat.Size() == 0 {
		t.Error("Log file is empty, expected content")
	}
}

func TestLoggingTransportClose(t *testing.T) {
	tmpDir := t.TempDir()
	logFile, err := os.Create(filepath.Join(tmpDir, "test.log"))
	if err != nil {
		t.Fatal(err)
	}

	transport := &LoggingTransport{
		Transport: http.DefaultTransport,
		LogFile:   logFile,
	}

	err = transport.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestLoggingTransportCloseWithNilFile(t *testing.T) {
	transport := &LoggingTransport{
		Transport: http.DefaultTransport,
		LogFile:   nil,
	}

	err := transport.Close()
	if err != nil {
		t.Errorf("Close() with nil file error = %v, want nil", err)
	}
}

func TestIsTFXLogEnabled(t *testing.T) {
	// This test verifies the function doesn't panic
	// and returns a boolean value
	result := IsTFXLogEnabled()
	t.Logf("IsTFXLogEnabled() = %v", result)
}

func TestNewHTTPClientWithLogging(t *testing.T) {
	t.Run("creates client", func(t *testing.T) {
		client, closer, err := NewHTTPClientWithLogging()
		if err != nil {
			t.Fatalf("NewHTTPClientWithLogging() error = %v", err)
		}
		if client == nil {
			t.Error("Expected non-nil client")
		}
		if closer != nil {
			defer closer.Close()
		}
	})
}
