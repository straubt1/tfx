package client

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LoggingTransport wraps an http.RoundTripper to log requests and responses to a file
type LoggingTransport struct {
	Transport http.RoundTripper
	LogFile   *os.File
	// If true, also write brief or debug output to stderr
	LogToTerminal bool
	// raw value of TFX_LOG (e.g. "debug")
	LogLevel string
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
		// still may want to log to terminal
		if !t.LogToTerminal {
			return
		}
	}

	timestamp := time.Now().Format(time.RFC3339)
	fmt.Fprintf(t.LogFile, "================================================================================\n")
	fmt.Fprintf(t.LogFile, "REQUEST @ %s\n", timestamp)
	fmt.Fprintf(t.LogFile, "================================================================================\n")

	// Dump the request with body
	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		if t.LogFile != nil {
			fmt.Fprintf(t.LogFile, "Error dumping request: %v\n", err)
		}
		if t.LogToTerminal {
			fmt.Fprintf(os.Stderr, "Error dumping request: %v\n", err)
		}
		return
	}

	if t.LogFile != nil {
		t.LogFile.Write(reqDump)
		t.LogFile.WriteString("\n")
	}

	if t.LogToTerminal {
		// If debug level, print full dump; otherwise print summary
		if strings.ToLower(t.LogLevel) == "debug" {
			fmt.Fprintln(os.Stderr, string(reqDump))
		} else {
			fmt.Fprintf(os.Stderr, "REQUEST %s %s\n", req.Method, req.URL.String())
		}
	}
}

func (t *LoggingTransport) logResponse(resp *http.Response) {
	if t.LogFile == nil && !t.LogToTerminal {
		return
	}

	timestamp := time.Now().Format(time.RFC3339)
	fmt.Fprintf(t.LogFile, "--------------------------------------------------------------------------------\n")
	fmt.Fprintf(t.LogFile, "RESPONSE @ %s\n", timestamp)
	fmt.Fprintf(t.LogFile, "--------------------------------------------------------------------------------\n")

	// Dump the response with body
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		if t.LogFile != nil {
			fmt.Fprintf(t.LogFile, "Error dumping response: %v\n", err)
		}
		if t.LogToTerminal {
			fmt.Fprintf(os.Stderr, "Error dumping response: %v\n", err)
		}
		return
	}

	if t.LogFile != nil {
		t.LogFile.Write(respDump)
		t.LogFile.WriteString("\n")
	}

	if t.LogToTerminal {
		if strings.ToLower(t.LogLevel) == "debug" {
			fmt.Fprintln(os.Stderr, string(respDump))
		} else {
			fmt.Fprintf(os.Stderr, "RESPONSE %s %s - Status: %s\n", resp.Request.Method, resp.Request.URL.String(), resp.Status)
		}
	}
}

func (t *LoggingTransport) logError(err error) {
	timestamp := time.Now().Format(time.RFC3339)
	if t.LogFile != nil {
		fmt.Fprintf(t.LogFile, "\n*** ERROR @ %s ***\n", timestamp)
		fmt.Fprintf(t.LogFile, "%v\n", err)
	}
	if t.LogToTerminal {
		fmt.Fprintf(os.Stderr, "*** ERROR @ %s ***\n%v\n", timestamp, err)
	}
}

// Close closes the log file
func (t *LoggingTransport) Close() error {
	if t.LogFile != nil {
		return t.LogFile.Close()
	}
	return nil
}

func IsTFXLogEnabled() bool {
	// Enabled when either TFX_LOG (terminal logging) or TFX_LOG_PATH (file logging) is set.
	// The actual env lookup happens at package init and values are cached.
	return tfxLogLevel != "" || tfxLogPath != ""
}

// NewHTTPClientWithLogging creates an HTTP client that logs all requests and responses to a file
// package-level cached env values
var (
	tfxLogLevel string
	tfxLogPath  string
)

func init() {
	tfxLogLevel = strings.TrimSpace(os.Getenv("TFX_LOG"))
	tfxLogPath = strings.TrimSpace(os.Getenv("TFX_LOG_PATH"))
}

// NewHTTPClientWithLogging creates an HTTP client that logs all requests and responses to a file
// if TFX_LOG_PATH is set, and/or to the terminal if TFX_LOG is set.
func NewHTTPClientWithLogging() (*http.Client, io.Closer, error) {
	var logFile *os.File
	var err error

	if tfxLogPath != "" {
		// Ensure directory exists
		if err = os.MkdirAll(tfxLogPath, 0755); err != nil {
			return nil, nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// Create a timestamped log file
		filename := fmt.Sprintf("tfx_http_%s_%d.log", time.Now().Format("20060102_150405"), os.Getpid())
		logFilePath := filepath.Join(tfxLogPath, filename)

		logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to open log file: %w", err)
		}

		// Write header to log file
		timestamp := time.Now().Format(time.RFC3339)
		// fmt.Fprintf(logFile, "\n\n")
		fmt.Fprintf(logFile, "################################################################################\n")
		fmt.Fprintf(logFile, "# TFX HTTP LOG - Started at %s\n", timestamp)
		fmt.Fprintf(logFile, "################################################################################\n")
	}

	transport := &LoggingTransport{
		Transport:     http.DefaultTransport,
		LogFile:       logFile,
		LogToTerminal: tfxLogLevel != "",
		LogLevel:      tfxLogLevel,
	}

	client := &http.Client{
		Transport: transport,
	}

	return client, transport, nil
}
