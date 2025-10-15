// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package logger

import (
	"context"
	"log/slog"
	"os"
)

var (
	// global logger instance
	std *slog.Logger
	// current log level
	currentLevel slog.Level
	// log path from TFX_LOG_PATH environment variable
	logPath string
)

// Init initializes the global logger with the specified log level string and log path
func Init(logLevel string, logPathEnv string) {
	level := ParseLevel(logLevel)
	currentLevel = level
	logPath = logPathEnv

	opts := &slog.HandlerOptions{
		Level: currentLevel,
	}

	handler := NewColorHandler(os.Stderr, opts)
	std = slog.New(handler)

	if IsEnabled(slog.LevelDebug) {
		Debug("Logger initialized", "level", LevelString(level))
		if logPath != "" {
			Debug("Log path configured", "path", logPath)
		}
	}
}

// GetLevel returns the current log level
func GetLevel() slog.Level {
	return currentLevel
}

// GetLogPath returns the configured log path
func GetLogPath() string {
	return logPath
}

// IsEnabled returns true if the logger is enabled at the given level
func IsEnabled(level slog.Level) bool {
	return currentLevel <= level
}

// Trace logs a trace message (using Debug-1 level in slog)
func Trace(msg string, keysAndValues ...interface{}) {
	if std != nil {
		std.Log(context.Background(), LevelTrace, msg, keysAndValues...)
	}
}

// Debug logs a debug message
func Debug(msg string, keysAndValues ...interface{}) {
	if std != nil {
		std.Debug(msg, keysAndValues...)
	}
}

// Info logs an informational message
func Info(msg string, keysAndValues ...interface{}) {
	if std != nil {
		std.Info(msg, keysAndValues...)
	}
}

// Warn logs a warning message
func Warn(msg string, keysAndValues ...interface{}) {
	if std != nil {
		std.Warn(msg, keysAndValues...)
	}
}

// Error logs an error message
func Error(msg string, keysAndValues ...interface{}) {
	if std != nil {
		std.Error(msg, keysAndValues...)
	}
}
