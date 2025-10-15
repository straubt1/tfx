// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package logger

import (
	"log/slog"
	"strings"
)

// Define custom log levels as constants
const (
	// LevelTrace is more verbose than Debug
	LevelTrace = slog.LevelDebug - 1
	// LevelOff disables all logging
	LevelOff = slog.Level(100)
)

// ParseLevel converts a string to a slog.Level
func ParseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "trace":
		return LevelTrace
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "off", "":
		return LevelOff
	default:
		return LevelOff
	}
}

// LevelString returns a string representation of the log level
func LevelString(level slog.Level) string {
	switch {
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
	case level >= LevelOff:
		return "OFF"
	default:
		return "UNKNOWN"
	}
}
