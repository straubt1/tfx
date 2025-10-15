package client

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"time"

	"github.com/straubt1/tfx/logger"
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
	// Log to file if enabled
	if t.LogFile != nil {
		timestamp := time.Now().Format(time.RFC3339)
		fmt.Fprintf(t.LogFile, "================================================================================\n")
		fmt.Fprintf(t.LogFile, "REQUEST @ %s\n", timestamp)
		fmt.Fprintf(t.LogFile, "================================================================================\n")

		// Dump the request with body
		reqDump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			fmt.Fprintf(t.LogFile, "Error dumping request: %v\n", err)
		} else {
			t.LogFile.Write(reqDump)
		}
	}

	// Log to logger based on log level
	if logger.IsEnabled(logger.LevelTrace) {
		// At TRACE level, log the full request dump
		reqDump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			logger.Error("Failed to dump HTTP request", "error", err)
		} else {
			logger.Trace("HTTP Request (full dump)", "request", string(reqDump))
		}
	} else if logger.IsEnabled(slog.LevelDebug) {
		// At DEBUG level, log request summary
		logger.Debug("HTTP Request", "method", req.Method, "url", req.URL.String())
	}
}

func (t *LoggingTransport) logResponse(resp *http.Response) {
	// Log to file if enabled
	if t.LogFile != nil {
		timestamp := time.Now().Format(time.RFC3339)
		fmt.Fprintf(t.LogFile, "--------------------------------------------------------------------------------\n")
		fmt.Fprintf(t.LogFile, "RESPONSE @ %s\n", timestamp)
		fmt.Fprintf(t.LogFile, "--------------------------------------------------------------------------------\n")

		// Dump the response with body
		respDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Fprintf(t.LogFile, "Error dumping response: %v\n", err)
		} else {
			t.LogFile.Write(respDump)
		}
	}

	// Log to logger based on log level
	if logger.IsEnabled(logger.LevelTrace) {
		// At TRACE level, log the full response dump
		respDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			logger.Error("Failed to dump HTTP response", "error", err)
		} else {
			logger.Trace("HTTP Response (full dump)", "response", string(respDump))
		}
	} else if logger.IsEnabled(slog.LevelDebug) {
		// At DEBUG level, log response summary
		logger.Debug("HTTP Response",
			"method", resp.Request.Method,
			"url", resp.Request.URL.String(),
			"status", resp.Status,
			"statusCode", resp.StatusCode)
	}
}

func (t *LoggingTransport) logError(err error) {
	// Log to file if enabled
	if t.LogFile != nil {
		timestamp := time.Now().Format(time.RFC3339)
		fmt.Fprintf(t.LogFile, "\n*** ERROR @ %s ***\n", timestamp)
		fmt.Fprintf(t.LogFile, "%v\n", err)
	}

	// Log to logger
	logger.Error("HTTP transport error", "error", err)
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
	return logger.IsEnabled(slog.LevelInfo) || logger.GetLogPath() != ""
}

// NewHTTPClientWithLogging creates an HTTP client that logs all requests and responses to a file
// if TFX_LOG_PATH is set, and/or to the terminal if TFX_LOG is set.
func NewHTTPClientWithLogging() (*http.Client, io.Closer, error) {
	var logFile *os.File
	var err error

	logPath := logger.GetLogPath()
	if logPath != "" {
		// Ensure directory exists
		if err = os.MkdirAll(logPath, 0755); err != nil {
			return nil, nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// Create a timestamped log file
		filename := fmt.Sprintf("tfx_http_%s_%d.log", time.Now().Format("20060102_150405"), os.Getpid())
		logFilePath := filepath.Join(logPath, filename)

		logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to open log file: %w", err)
		}

		// Write header to log file
		timestamp := time.Now().Format(time.RFC3339)
		fmt.Fprintf(logFile, "################################################################################\n")
		fmt.Fprintf(logFile, "# TFX HTTP LOG - Started at %s\n", timestamp)
		fmt.Fprintf(logFile, "################################################################################\n")

		logger.Info("HTTP logging to file enabled", "path", logFilePath)
	}

	transport := &LoggingTransport{
		Transport: http.DefaultTransport,
		LogFile:   logFile,
	}

	client := &http.Client{
		Transport: transport,
	}

	return client, transport, nil
}
