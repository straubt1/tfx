// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/straubt1/tfx/output"
)

// Regex patterns to redact sensitive information from HTTP dumps
var (
	// Matches Authorization header with any token/credentials
	authHeaderRegex = regexp.MustCompile(`(?i)(Authorization:\s+)([^\r\n]+)`)
	// Matches other potential sensitive headers
	cookieHeaderRegex = regexp.MustCompile(`(?i)(Cookie:\s+)([^\r\n]+)`)
	apiKeyHeaderRegex = regexp.MustCompile(`(?i)(X-Api-Key:\s+)([^\r\n]+)`)
)

// redactSensitiveData replaces sensitive information in HTTP dumps with [REDACTED]
func redactSensitiveData(dump []byte) []byte {
	result := dump

	// Redact Authorization header
	result = authHeaderRegex.ReplaceAll(result, []byte("${1}[REDACTED]"))

	// Redact Cookie header
	result = cookieHeaderRegex.ReplaceAll(result, []byte("${1}[REDACTED]"))

	// Redact API Key header
	result = apiKeyHeaderRegex.ReplaceAll(result, []byte("${1}[REDACTED]"))

	return result
}

// LoggingTransport wraps an http.RoundTripper to log requests and responses.
// It optionally publishes APIEvents to an APIEventBus for the TUI inspector panel.
type LoggingTransport struct {
	Transport http.RoundTripper
	LogFile   *os.File
	EventBus  *APIEventBus // nil when not running in TUI mode
}

// RoundTrip implements the http.RoundTripper interface with logging and event publishing.
func (t *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	// Capture the request body for the event bus before logging alters req.Body.
	// This is only done for methods that carry a body (POST, PATCH, PUT) and only
	// when the event bus is wired up. DumpRequestOut replaces req.Body with a
	// re-readable buffer so subsequent reads (logRequest, transport) still work.
	var reqBodyStr string
	if t.EventBus != nil {
		reqBodyStr = captureReqBody(req)
	}

	// Log request (may also DumpRequestOut internally for TRACE; safe after our capture).
	t.logRequest(req)

	// Perform the actual request.
	resp, err := t.Transport.RoundTrip(req)
	duration := time.Since(start)

	if err != nil {
		t.logError(err)
		if t.EventBus != nil {
			t.EventBus.Send(APIEvent{
				Timestamp:  time.Now(),
				Method:     req.Method,
				URL:        req.URL.String(),
				Path:       req.URL.Path,
				Duration:   duration,
				ReqHeaders: formatHeaders(req.Header),
				ReqBody:    reqBodyStr,
				Err:        err.Error(),
			})
		}
		return nil, err
	}

	// Capture response dump ONCE — shared by file logger, trace logger, and event bus.
	// DumpResponse(resp, true) reads and replaces resp.Body with a re-readable buffer.
	var respDump []byte
	needsDump := t.LogFile != nil ||
		output.Get().Logger().IsEnabled(output.LevelTrace) ||
		t.EventBus != nil
	if needsDump {
		var dumpErr error
		respDump, dumpErr = httputil.DumpResponse(resp, true)
		if dumpErr != nil {
			respDump = nil // non-fatal — logging/inspector miss this response
		}
	}

	// Log response (file + terminal) using the pre-captured dump.
	t.logResponseFromDump(resp, respDump)

	// Publish to the TUI event bus.
	if t.EventBus != nil {
		t.EventBus.Send(APIEvent{
			Timestamp:   time.Now(),
			Method:      req.Method,
			URL:         req.URL.String(),
			Path:        pathFromURL(req.URL),
			StatusCode:  resp.StatusCode,
			Duration:    duration,
			ReqHeaders:  formatHeaders(req.Header),
			ReqBody:     reqBodyStr,
			RespHeaders: formatHeaders(resp.Header),
			RespBody:    extractAndPrettyBody(respDump),
		})
	}

	return resp, nil
}

// formatHeaders formats http.Header as sorted "Name: value" strings, redacting
// known sensitive headers (Authorization, Cookie, X-Api-Key).
func formatHeaders(h http.Header) []string {
	sensitive := map[string]bool{
		"Authorization": true,
		"Cookie":        true,
		"X-Api-Key":     true,
	}

	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make([]string, 0, len(keys))
	for _, k := range keys {
		val := strings.Join(h[k], ", ")
		if sensitive[http.CanonicalHeaderKey(k)] {
			val = "[REDACTED]"
		}
		out = append(out, k+": "+val)
	}
	return out
}

// captureReqBody reads the request body (for POST/PATCH/PUT) and returns it as a
// string. It replaces req.Body with a fresh re-readable buffer so the actual
// transport can still read the original bytes.
func captureReqBody(req *http.Request) string {
	if req.Body == nil || req.Body == http.NoBody {
		return ""
	}
	// Only capture for methods that carry bodies.
	switch req.Method {
	case http.MethodPost, http.MethodPatch, http.MethodPut:
	default:
		return ""
	}
	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return ""
	}
	// Extract body portion (everything after the header/blank-line separator).
	parts := bytes.SplitN(dump, []byte("\r\n\r\n"), 2)
	if len(parts) < 2 {
		return ""
	}
	return prettyJSONBytes(parts[1])
}

// pathFromURL returns just the path (and query) portion of u, suitable for display.
func pathFromURL(u *url.URL) string {
	p := u.Path
	if u.RawQuery != "" {
		p += "?" + u.RawQuery
	}
	return p
}

// extractAndPrettyBody extracts the HTTP body from a full response dump and
// pretty-prints it if it is valid JSON. Returns an empty string on failure.
func extractAndPrettyBody(dump []byte) string {
	if len(dump) == 0 {
		return ""
	}
	parts := bytes.SplitN(dump, []byte("\r\n\r\n"), 2)
	if len(parts) < 2 {
		return ""
	}
	return prettyJSONBytes(parts[1])
}

// prettyJSONBytes attempts to pretty-print b as JSON. Returns the original bytes
// as a string when b is not valid JSON (e.g., binary or plain text).
func prettyJSONBytes(b []byte) string {
	b = bytes.TrimSpace(b)
	if len(b) == 0 {
		return ""
	}
	var buf bytes.Buffer
	if err := json.Indent(&buf, b, "", "  "); err == nil {
		return buf.String()
	}
	return string(b)
}

func (t *LoggingTransport) logRequest(req *http.Request) {
	// Log to file if enabled
	if t.LogFile != nil {
		timestamp := time.Now().Format(time.RFC3339)
		_, _ = fmt.Fprintf(t.LogFile, "================================================================================\n")
		_, _ = fmt.Fprintf(t.LogFile, "REQUEST @ %s\n", timestamp)
		_, _ = fmt.Fprintf(t.LogFile, "================================================================================\n")

		// Dump the request with body
		reqDump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			_, _ = fmt.Fprintf(t.LogFile, "Error dumping request: %v\n", err)
		} else {
			// Redact sensitive data before writing to file
			redactedDump := redactSensitiveData(reqDump)
			_, _ = t.LogFile.Write(redactedDump)
		}
	}

	// Log to logger based on log level
	if output.Get().Logger().IsEnabled(output.LevelTrace) {
		// At TRACE level, log the full request dump
		reqDump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			output.Get().Logger().Error("Failed to dump HTTP request", "error", err)
		} else {
			// Redact sensitive data before logging
			redactedDump := redactSensitiveData(reqDump)
			output.Get().Logger().Trace("HTTP Request (full dump)", "request", string(redactedDump))
		}
	} else if output.Get().Logger().IsEnabled(slog.LevelDebug) {
		// At DEBUG level, log request summary
		output.Get().Logger().Debug("HTTP Request", "method", req.Method, "url", req.URL.String())
	}
}

// logResponseFromDump logs the response using the pre-captured dump bytes.
// The dump was already obtained by the caller via DumpResponse(resp, true).
func (t *LoggingTransport) logResponseFromDump(resp *http.Response, dump []byte) {
	// Log to file if enabled
	if t.LogFile != nil {
		timestamp := time.Now().Format(time.RFC3339)
		_, _ = fmt.Fprintf(t.LogFile, "--------------------------------------------------------------------------------\n")
		_, _ = fmt.Fprintf(t.LogFile, "RESPONSE @ %s\n", timestamp)
		_, _ = fmt.Fprintf(t.LogFile, "--------------------------------------------------------------------------------\n")

		if len(dump) > 0 {
			redactedDump := redactSensitiveData(dump)
			_, _ = t.LogFile.Write(redactedDump)
		}
	}

	// Log to logger based on log level
	if output.Get().Logger().IsEnabled(output.LevelTrace) {
		if len(dump) > 0 {
			redactedDump := redactSensitiveData(dump)
			output.Get().Logger().Trace("HTTP Response (full dump)", "response", string(redactedDump))
		}
	} else if output.Get().Logger().IsEnabled(slog.LevelDebug) {
		output.Get().Logger().Debug("HTTP Response",
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
	output.Get().Logger().Error("HTTP transport error", "error", err)
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
	return output.Get().Logger().IsEnabled(slog.LevelInfo) || output.Get().Logger().GetLogPath() != ""
}

// NewHTTPClientWithLogging creates an HTTP client that logs all requests and responses to a file
// if TFX_LOG_PATH is set, and/or to the terminal if TFX_LOG is set.
// An optional EventBus may be provided to also publish events for the TUI inspector panel.
func NewHTTPClientWithLogging(bus *APIEventBus) (*http.Client, io.Closer, error) {
	var logFile *os.File
	var err error

	logPath := output.Get().Logger().GetLogPath()
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

		output.Get().Logger().Info("HTTP logging to file enabled", "path", logFilePath)
	}

	transport := &LoggingTransport{
		Transport: http.DefaultTransport,
		LogFile:   logFile,
		EventBus:  bus,
	}

	client := &http.Client{
		Transport: transport,
	}

	return client, transport, nil
}

// methodColor returns an ANSI color prefix for a given HTTP method for display purposes.
// Not used internally — exported so the TUI debugpanel can reference it.
func MethodDisplayLabel(method string) string {
	switch strings.ToUpper(method) {
	case "GET":
		return "GET   "
	case "POST":
		return "POST  "
	case "PATCH":
		return "PATCH "
	case "PUT":
		return "PUT   "
	case "DELETE":
		return "DELETE"
	default:
		return method
	}
}
