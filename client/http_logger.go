package client

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

// LoggingTransport wraps an http.RoundTripper to log requests and responses to a file
type LoggingTransport struct {
	Transport http.RoundTripper
	LogFile   *os.File
}

// RoundTrip implements the http.RoundTripper interface with logging
func (t *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Log request
	t.logRequest(req)

	// Perform the actual request
	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		t.logError(err)
		return nil, err
	}

	// Log response
	t.logResponse(resp)

	return resp, nil
}

func (t *LoggingTransport) logRequest(req *http.Request) {
	if t.LogFile == nil {
		return
	}

	timestamp := time.Now().Format(time.RFC3339)
	fmt.Fprintf(t.LogFile, "\n================================================================================\n")
	fmt.Fprintf(t.LogFile, "REQUEST @ %s\n", timestamp)
	fmt.Fprintf(t.LogFile, "================================================================================\n")

	// Dump the request with body
	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		fmt.Fprintf(t.LogFile, "Error dumping request: %v\n", err)
		return
	}

	t.LogFile.Write(reqDump)
	t.LogFile.WriteString("\n")
}

func (t *LoggingTransport) logResponse(resp *http.Response) {
	if t.LogFile == nil {
		return
	}

	timestamp := time.Now().Format(time.RFC3339)
	fmt.Fprintf(t.LogFile, "\n--------------------------------------------------------------------------------\n")
	fmt.Fprintf(t.LogFile, "RESPONSE @ %s\n", timestamp)
	fmt.Fprintf(t.LogFile, "--------------------------------------------------------------------------------\n")

	// Dump the response with body
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Fprintf(t.LogFile, "Error dumping response: %v\n", err)
		return
	}

	t.LogFile.Write(respDump)
	t.LogFile.WriteString("\n")
}

func (t *LoggingTransport) logError(err error) {
	if t.LogFile == nil {
		return
	}

	timestamp := time.Now().Format(time.RFC3339)
	fmt.Fprintf(t.LogFile, "\n*** ERROR @ %s ***\n", timestamp)
	fmt.Fprintf(t.LogFile, "%v\n", err)
}

// Close closes the log file
func (t *LoggingTransport) Close() error {
	if t.LogFile != nil {
		return t.LogFile.Close()
	}
	return nil
}

// NewHTTPClientWithLogging creates an HTTP client that logs all requests and responses to a file
func NewHTTPClientWithLogging(logFilePath string) (*http.Client, io.Closer, error) {
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Write header to log file
	timestamp := time.Now().Format(time.RFC3339)
	fmt.Fprintf(logFile, "\n\n")
	fmt.Fprintf(logFile, "################################################################################\n")
	fmt.Fprintf(logFile, "# TFX HTTP LOG - Started at %s\n", timestamp)
	fmt.Fprintf(logFile, "################################################################################\n")

	transport := &LoggingTransport{
		Transport: http.DefaultTransport,
		LogFile:   logFile,
	}

	client := &http.Client{
		Transport: transport,
	}

	return client, transport, nil
}
