// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package output

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/logrusorgru/aurora"
)

// Custom log level for trace
const LevelTrace = slog.Level(-8)

// Logger wraps slog for structured logging
type Logger struct {
	slog  *slog.Logger
	level slog.Level
	path  string
}

// NewLogger creates a new logger instance
func NewLogger(logLevel string, logPath string) *Logger {
	level := parseLevel(logLevel)

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := newColorHandler(os.Stderr, opts)
	logger := slog.New(handler)

	l := &Logger{
		slog:  logger,
		level: level,
		path:  logPath,
	}

	if l.IsEnabled(slog.LevelDebug) {
		l.Debug("Logger initialized", "level", levelString(level))
		if logPath != "" {
			l.Debug("Log path configured", "path", logPath)
		}
	}

	return l
}

// parseLevel parses a log level string
func parseLevel(levelStr string) slog.Level {
	switch levelStr {
	case "TRACE", "trace":
		return LevelTrace
	case "DEBUG", "debug":
		return slog.LevelDebug
	case "INFO", "info":
		return slog.LevelInfo
	case "WARN", "warn":
		return slog.LevelWarn
	case "ERROR", "error":
		return slog.LevelError
	case "NONE", "none", "":
		// Disable logging by setting a very high level
		return slog.Level(1000)
	default:
		// Default to none (disabled) when unrecognized
		return slog.Level(1000)
	}
}

// levelString returns the string representation of a log level
func levelString(level slog.Level) string {
	switch {
	case level >= slog.Level(1000):
		return "NONE"
	case level == LevelTrace:
		return "TRACE"
	case level == slog.LevelDebug:
		return "DEBUG"
	case level == slog.LevelInfo:
		return "INFO"
	case level == slog.LevelWarn:
		return "WARN"
	case level == slog.LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// GetLevel returns the current log level
func (l *Logger) GetLevel() slog.Level {
	return l.level
}

// GetLogPath returns the configured log path
func (l *Logger) GetLogPath() string {
	return l.path
}

// IsEnabled returns true if the logger is enabled at the given level
func (l *Logger) IsEnabled(level slog.Level) bool {
	return l.level <= level
}

// Trace logs a trace message
func (l *Logger) Trace(msg string, keysAndValues ...interface{}) {
	l.slog.Log(context.Background(), LevelTrace, msg, keysAndValues...)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	l.slog.Debug(msg, keysAndValues...)
}

// Info logs an informational message
func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.slog.Info(msg, keysAndValues...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	l.slog.Warn(msg, keysAndValues...)
}

// Error logs an error message
func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	l.slog.Error(msg, keysAndValues...)
}

// colorHandler is a custom slog.Handler that formats logs with colors
type colorHandler struct {
	opts slog.HandlerOptions
	out  io.Writer
}

// newColorHandler creates a new colorHandler
func newColorHandler(out io.Writer, opts *slog.HandlerOptions) *colorHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &colorHandler{
		opts: *opts,
		out:  out,
	}
}

// Enabled reports whether the handler handles records at the given level
func (h *colorHandler) Enabled(_ context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

// Handle formats and writes the log record
func (h *colorHandler) Handle(_ context.Context, r slog.Record) error {
	timestamp := r.Time.Format("15:04:05")
	level := h.formatLevel(r.Level)

	msg := fmt.Sprintf("%s %s %s",
		aurora.Faint(fmt.Sprintf("[%s]", timestamp)),
		level,
		r.Message,
	)

	// Add attributes
	r.Attrs(func(a slog.Attr) bool {
		msg += fmt.Sprintf(" %s=%s",
			aurora.Faint(a.Key),
			aurora.Faint(a.Value.String()),
		)
		return true
	})

	fmt.Fprintln(h.out, msg)
	return nil
}

// WithAttrs returns a new Handler with additional attributes
func (h *colorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

// WithGroup returns a new Handler with a group name
func (h *colorHandler) WithGroup(name string) slog.Handler {
	return h
}

// formatLevel formats the log level with color
func (h *colorHandler) formatLevel(level slog.Level) string {
	var levelStr string
	var colored string

	switch {
	case level < slog.LevelDebug:
		levelStr = "TRACE"
		colored = aurora.Faint(fmt.Sprintf("[%s]", levelStr)).String()
	case level < slog.LevelInfo:
		levelStr = "DEBUG"
		colored = aurora.Cyan(fmt.Sprintf("[%s]", levelStr)).String()
	case level < slog.LevelWarn:
		levelStr = "INFO"
		colored = aurora.Green(fmt.Sprintf("[%s]", levelStr)).String()
	case level < slog.LevelError:
		levelStr = "WARN"
		colored = aurora.Yellow(fmt.Sprintf("[%s]", levelStr)).String()
	default:
		levelStr = "ERROR"
		colored = aurora.Red(fmt.Sprintf("[%s]", levelStr)).String()
	}

	return colored
}
